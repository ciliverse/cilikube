package service

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/ciliverse/cilikube/pkg/k8s"
)

// --- GatewayClassClient (Cluster-scoped) ---
type GatewayClassClient struct{}

func (c *GatewayClassClient) Get(ctx context.Context, clientset kubernetes.Interface, _, name string, opts metav1.GetOptions) (*gatewayv1.GatewayClass, error) {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return nil, err
	}
	return gw.GatewayV1().GatewayClasses().Get(ctx, name, opts)
}
func (c *GatewayClassClient) List(ctx context.Context, clientset kubernetes.Interface, _ string, opts metav1.ListOptions) (runtime.Object, error) {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return nil, err
	}
	return gw.GatewayV1().GatewayClasses().List(ctx, opts)
}
func (c *GatewayClassClient) Create(ctx context.Context, clientset kubernetes.Interface, _ string, obj *gatewayv1.GatewayClass, opts metav1.CreateOptions) (*gatewayv1.GatewayClass, error) {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return nil, err
	}
	return gw.GatewayV1().GatewayClasses().Create(ctx, obj, opts)
}
func (c *GatewayClassClient) Update(ctx context.Context, clientset kubernetes.Interface, _ string, obj *gatewayv1.GatewayClass, opts metav1.UpdateOptions) (*gatewayv1.GatewayClass, error) {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return nil, err
	}
	return gw.GatewayV1().GatewayClasses().Update(ctx, obj, opts)
}
func (c *GatewayClassClient) Delete(ctx context.Context, clientset kubernetes.Interface, _, name string, opts metav1.DeleteOptions) error {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return err
	}
	return gw.GatewayV1().GatewayClasses().Delete(ctx, name, opts)
}
func (c *GatewayClassClient) Watch(ctx context.Context, clientset kubernetes.Interface, _ string, opts metav1.ListOptions) (watch.Interface, error) {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return nil, err
	}
	return gw.GatewayV1().GatewayClasses().Watch(ctx, opts)
}

// --- GatewayClient (Namespaced) ---
type GatewayClient struct{}

func (c *GatewayClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*gatewayv1.Gateway, error) {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return nil, err
	}
	return gw.GatewayV1().Gateways(namespace).Get(ctx, name, opts)
}
func (c *GatewayClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return nil, err
	}
	return gw.GatewayV1().Gateways(namespace).List(ctx, opts)
}
func (c *GatewayClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *gatewayv1.Gateway, opts metav1.CreateOptions) (*gatewayv1.Gateway, error) {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return nil, err
	}
	return gw.GatewayV1().Gateways(namespace).Create(ctx, obj, opts)
}
func (c *GatewayClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *gatewayv1.Gateway, opts metav1.UpdateOptions) (*gatewayv1.Gateway, error) {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return nil, err
	}
	return gw.GatewayV1().Gateways(namespace).Update(ctx, obj, opts)
}
func (c *GatewayClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return err
	}
	return gw.GatewayV1().Gateways(namespace).Delete(ctx, name, opts)
}
func (c *GatewayClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return nil, err
	}
	return gw.GatewayV1().Gateways(namespace).Watch(ctx, opts)
}

// --- HTTPRouteClient (Namespaced) ---
type HTTPRouteClient struct{}

func (c *HTTPRouteClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*gatewayv1.HTTPRoute, error) {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return nil, err
	}
	return gw.GatewayV1().HTTPRoutes(namespace).Get(ctx, name, opts)
}
func (c *HTTPRouteClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return nil, err
	}
	return gw.GatewayV1().HTTPRoutes(namespace).List(ctx, opts)
}
func (c *HTTPRouteClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *gatewayv1.HTTPRoute, opts metav1.CreateOptions) (*gatewayv1.HTTPRoute, error) {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return nil, err
	}
	return gw.GatewayV1().HTTPRoutes(namespace).Create(ctx, obj, opts)
}
func (c *HTTPRouteClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *gatewayv1.HTTPRoute, opts metav1.UpdateOptions) (*gatewayv1.HTTPRoute, error) {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return nil, err
	}
	return gw.GatewayV1().HTTPRoutes(namespace).Update(ctx, obj, opts)
}
func (c *HTTPRouteClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return err
	}
	return gw.GatewayV1().HTTPRoutes(namespace).Delete(ctx, name, opts)
}
func (c *HTTPRouteClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	gw, err := k8s.GatewayClientFor(clientset)
	if err != nil {
		return nil, err
	}
	return gw.GatewayV1().HTTPRoutes(namespace).Watch(ctx, opts)
}
