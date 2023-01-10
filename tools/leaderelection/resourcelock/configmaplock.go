/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resourcelock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/typed"
	corev1 "k8s.io/client-go/tools/leaderelection/resourcelock/internal/apis/core/v1"
)

// TODO: This is almost a exact replica of Endpoints lock.
// going forwards as we self host more and more components
// and use ConfigMaps as the means to pass that configuration
// data we will likely move to deprecate the Endpoints lock.

type configMapLock struct {
	// ConfigMapMeta should contain a Name and a Namespace of a
	// ConfigMapMeta object that the LeaderElector will attempt to lead.
	ConfigMapMeta metav1.ObjectMeta
	Client        dynamic.Interface
	LockConfig    ResourceLockConfig
	cm            *corev1.ConfigMap
}

func (ll *configMapLock) ConfigMaps() typed.NamespaceClient[corev1.ConfigMap] {
	return typed.NewTypedNamespaceScoped[corev1.ConfigMap](ll.Client, schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "configmaps",
	})
}

// Get returns the election record from a ConfigMap Annotation
func (cml *configMapLock) Get(ctx context.Context) (*LeaderElectionRecord, []byte, error) {
	var record LeaderElectionRecord
	cm, err := cml.ConfigMaps().Namespace(cml.ConfigMapMeta.Namespace).Get(ctx, cml.ConfigMapMeta.Name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}
	cml.cm = cm
	if cml.cm.Annotations == nil {
		cml.cm.Annotations = make(map[string]string)
	}
	recordStr, found := cml.cm.Annotations[LeaderElectionRecordAnnotationKey]
	recordBytes := []byte(recordStr)
	if found {
		if err := json.Unmarshal(recordBytes, &record); err != nil {
			return nil, nil, err
		}
	}
	return &record, recordBytes, nil
}

// Create attempts to create a LeaderElectionRecord annotation
func (cml *configMapLock) Create(ctx context.Context, ler LeaderElectionRecord) error {
	recordBytes, err := json.Marshal(ler)
	if err != nil {
		return err
	}
	cml.cm, err = cml.ConfigMaps().Namespace(cml.ConfigMapMeta.Namespace).Create(ctx, &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cml.ConfigMapMeta.Name,
			Namespace: cml.ConfigMapMeta.Namespace,
			Annotations: map[string]string{
				LeaderElectionRecordAnnotationKey: string(recordBytes),
			},
		},
	}, metav1.CreateOptions{})
	return err
}

// Update will update an existing annotation on a given resource.
func (cml *configMapLock) Update(ctx context.Context, ler LeaderElectionRecord) error {
	if cml.cm == nil {
		return errors.New("configmap not initialized, call get or create first")
	}
	recordBytes, err := json.Marshal(ler)
	if err != nil {
		return err
	}
	if cml.cm.Annotations == nil {
		cml.cm.Annotations = make(map[string]string)
	}
	cml.cm.Annotations[LeaderElectionRecordAnnotationKey] = string(recordBytes)
	cm, err := cml.ConfigMaps().Namespace(cml.ConfigMapMeta.Namespace).Update(ctx, cml.cm, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	cml.cm = cm
	return nil
}

// RecordEvent in leader election while adding meta-data
func (cml *configMapLock) RecordEvent(s string) {
	if cml.LockConfig.EventRecorder == nil {
		return
	}
	events := fmt.Sprintf("%v %v", cml.LockConfig.Identity, s)
	subject := &corev1.ConfigMap{ObjectMeta: cml.cm.ObjectMeta}
	// Populate the type meta, so we don't have to get it from the schema
	subject.Kind = "ConfigMap"
	subject.APIVersion = "v1" //internalapi.SchemeGroupVersion.String()
	cml.LockConfig.EventRecorder.Eventf(subject, corev1.EventTypeNormal, "LeaderElection", events)
}

// Describe is used to convert details on current resource lock
// into a string
func (cml *configMapLock) Describe() string {
	return fmt.Sprintf("%v/%v", cml.ConfigMapMeta.Namespace, cml.ConfigMapMeta.Name)
}

// Identity returns the Identity of the lock
func (cml *configMapLock) Identity() string {
	return cml.LockConfig.Identity
}
