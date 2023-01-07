package typed

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
)

func NewTypedNamespaceScoped[T any](client dynamic.Interface, gvr schema.GroupVersionResource) NamespaceClient[T] {
	panic("NewTypedNamespaceScoped not implemented")
}

func NewTypedClusterScoped[T any](client dynamic.Interface, gvr schema.GroupVersionResource) Client[T] {
	panic("NewTypedClusterScoped not implemented")
}

type Client[T any] interface {
	Create(ctx context.Context, obj *T, options metav1.CreateOptions, subresources ...string) (*T, error)
	Update(ctx context.Context, obj *T, options metav1.UpdateOptions, subresources ...string) (*T, error)
	UpdateStatus(ctx context.Context, obj *T, options metav1.UpdateOptions) (*T, error)
	Delete(ctx context.Context, name string, options metav1.DeleteOptions, subresources ...string) error
	DeleteCollection(ctx context.Context, options metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(ctx context.Context, name string, options metav1.GetOptions, subresources ...string) (*T, error)
	// List(ctx context.Context, opts metav1.ListOptions) (*TList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, options metav1.PatchOptions, subresources ...string) (*T, error)
	Apply(ctx context.Context, name string, obj *T, options metav1.ApplyOptions, subresources ...string) (*T, error)
	ApplyStatus(ctx context.Context, name string, obj *T, options metav1.ApplyOptions) (*T, error)
}

type NamespaceClient[T any] interface {
	Namespace(string) Client[T]

	// TODO: List / Watch etc
}
