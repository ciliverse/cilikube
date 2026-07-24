package k8s

import (
	"context"
	"os"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestShowcaseClientListsSeededResources(t *testing.T) {
	c := NewShowcaseClient(ShowcaseSeedObjects()...)
	ctx := context.Background()

	ns, err := c.Clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(ns.Items) < 5 {
		t.Fatalf("expected rich namespaces, got %d", len(ns.Items))
	}

	pods, err := c.Clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(pods.Items) < 15 {
		t.Fatalf("expected many pods, got %d", len(pods.Items))
	}

	deps, err := c.Clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(deps.Items) < 5 {
		t.Fatalf("expected deployments, got %d", len(deps.Items))
	}

	events, err := c.Clientset.CoreV1().Events("").List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(events.Items) < 10 {
		t.Fatalf("expected events, got %d", len(events.Items))
	}
	for _, p := range pods.Items {
		if p.CreationTimestamp.IsZero() {
			t.Fatalf("pod %s missing CreationTimestamp", p.Name)
		}
	}

	gw, err := GatewayClientFor(c.Clientset)
	if err != nil {
		t.Fatal(err)
	}
	classes, err := gw.GatewayV1().GatewayClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(classes.Items) < 1 {
		t.Fatalf("expected gateway classes, got %d", len(classes.Items))
	}
	routes, err := gw.GatewayV1().HTTPRoutes("").List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(routes.Items) < 1 {
		t.Fatalf("expected httproutes, got %d", len(routes.Items))
	}

	if !IsShowcaseConfig(c.Config) {
		t.Fatal("expected showcase rest config host")
	}
}

func TestIsShowcaseEnv(t *testing.T) {
	t.Setenv("CILIKUBE_SHOWCASE", "")
	t.Setenv("CILIKUBE_MODE", "")
	if IsShowcase() {
		t.Fatal("default must be false")
	}
	t.Setenv("CILIKUBE_SHOWCASE", "1")
	if !IsShowcase() {
		t.Fatal("SHOWCASE=1 should enable")
	}
	_ = os.Unsetenv("CILIKUBE_SHOWCASE")
}
