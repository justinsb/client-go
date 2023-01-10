package typed

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type namespacedClient[T any] struct {
	NamespaceClient[T]
	client dynamic.Interface
	gvr    schema.GroupVersionResource
}
