package service

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestTimelineAppGroup(t *testing.T) {
	if TimelineAppGroup(map[string]string{"app": "web"}) != "web" {
		t.Fatal("expected app label")
	}
	if TimelineAppGroup(nil) != "_ungrouped" {
		t.Fatal("expected ungrouped")
	}
}

func TestClassifyDeploymentStatus(t *testing.T) {
	one := int32(1)
	d := &appsv1.Deployment{
		Spec:   appsv1.DeploymentSpec{Replicas: &one},
		Status: appsv1.DeploymentStatus{ReadyReplicas: 1, UpdatedReplicas: 1},
	}
	st, _ := ClassifyDeploymentStatus(d)
	if st != TLHealthy {
		t.Fatalf("want healthy got %s", st)
	}
	d.Status.ReadyReplicas = 0
	d.Status.UnavailableReplicas = 1
	st, _ = ClassifyDeploymentStatus(d)
	if st != TLUnhealthy {
		t.Fatalf("want unhealthy got %s", st)
	}
}

func TestClassifyPodStatus(t *testing.T) {
	p := &corev1.Pod{
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
			ContainerStatuses: []corev1.ContainerStatus{
				{Ready: true, State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}},
			},
		},
	}
	st, _ := ClassifyPodStatus(p)
	if st != TLHealthy {
		t.Fatalf("want healthy got %s", st)
	}
	p.Status.ContainerStatuses[0].Ready = false
	p.Status.ContainerStatuses[0].State = corev1.ContainerState{
		Waiting: &corev1.ContainerStateWaiting{Reason: "CrashLoopBackOff"},
	}
	st, _ = ClassifyPodStatus(p)
	if st != TLUnhealthy {
		t.Fatalf("want unhealthy got %s", st)
	}
}

func TestMapEventMarker(t *testing.T) {
	if MapEventMarker("Warning", "BackOff") != "warning" {
		t.Fatal("warning")
	}
	if MapEventMarker("Normal", "Created") != "created" {
		t.Fatal("created")
	}
	if MapEventMarker("Normal", "Killing") != "deleted" {
		t.Fatal("deleted")
	}
}

func TestTimelineHref(t *testing.T) {
	if TimelineHref("Deployment", "ns", "x") != "/deployments/ns/x" {
		t.Fatal(TimelineHref("Deployment", "ns", "x"))
	}
}
