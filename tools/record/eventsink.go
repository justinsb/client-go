package record

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/typed"
	v1 "k8s.io/client-go/tools/record/internal/apis/core/v1"
)

type eventSink struct {
	client typed.NamespaceClient[v1.Event]
}

func NewEventSink(client dynamic.Interface) EventSink {
	typedClient := typed.NewTypedNamespaceScoped[v1.Event](client, schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "events",
	})

	return &eventSink{client: typedClient}
}
func (s *eventSink) Create(event *v1.Event) (*v1.Event, error) {
	ctx := context.TODO()
	return s.client.Namespace(event.Namespace).Create(ctx, event, metav1.CreateOptions{})
}

func (s *eventSink) Update(event *v1.Event) (*v1.Event, error) {
	ctx := context.TODO()
	return s.client.Namespace(event.Namespace).Update(ctx, event, metav1.UpdateOptions{})

}

func (s *eventSink) Patch(oldEvent *v1.Event, data []byte) (*v1.Event, error) {
	ctx := context.TODO()
	return s.client.Namespace(oldEvent.Namespace).Patch(ctx, oldEvent.Name, types.StrategicMergePatchType, data, metav1.PatchOptions{})
}
