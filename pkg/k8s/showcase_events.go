package k8s

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func ago(d time.Duration) metav1.Time {
	return metav1.NewTime(time.Now().Add(-d))
}

func event(ns, name, kind, objName, etype, reason, message string, first, last time.Duration, count int32) *corev1.Event {
	return &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Namespace:         ns,
			CreationTimestamp: ago(first),
		},
		InvolvedObject: corev1.ObjectReference{
			Kind:      kind,
			Namespace: ns,
			Name:      objName,
		},
		Reason:         reason,
		Message:        message,
		Type:           etype,
		Count:          count,
		FirstTimestamp: ago(first),
		LastTimestamp:  ago(last),
		Source:         corev1.EventSource{Component: "showcase-controller", Host: "demo-worker-1"},
	}
}

// ShowcaseEventObjects returns a dense stream of Normal/Warning events for the UI.
func ShowcaseEventObjects() []runtime.Object {
	evs := []*corev1.Event{
		event("default", "web-frontend-7d9f8b-abc12.started", "Pod", "web-frontend-7d9f8b-abc12",
			corev1.EventTypeNormal, "Started", "Container web-frontend started (showcase)", 2*time.Hour, 2*time.Hour, 1),
		event("default", "web-frontend-7d9f8b-abc12.pulled", "Pod", "web-frontend-7d9f8b-abc12",
			corev1.EventTypeNormal, "Pulled", "Successfully pulled image ghcr.io/ciliverse/web-frontend:showcase", 2*time.Hour+time.Minute, 2*time.Hour+time.Minute, 1),
		event("default", "web-frontend.scaled", "Deployment", "web-frontend",
			corev1.EventTypeNormal, "ScalingReplicaSet", "Scaled up replica set web-frontend-7d9f8b to 3", 6*time.Hour, 3*time.Hour, 2),
		event("default", "api-gateway-6c4d5-jkl78.unhealthy", "Pod", "api-gateway-6c4d5-jkl78",
			corev1.EventTypeWarning, "Unhealthy", "Readiness probe failed: connection refused (simulated)", 45*time.Minute, 12*time.Minute, 4),
		event("default", "api-gateway-6c4d5-jkl78.recovered", "Pod", "api-gateway-6c4d5-jkl78",
			corev1.EventTypeNormal, "Created", "Created container api-gateway after probe recovery", 10*time.Minute, 10*time.Minute, 1),
		event("production", "orders-api-8f2a1-aa111.scheduled", "Pod", "orders-api-8f2a1-aa111",
			corev1.EventTypeNormal, "Scheduled", "Successfully assigned production/orders-api-8f2a1-aa111 to demo-worker-1", 8*time.Hour, 8*time.Hour, 1),
		event("production", "orders-api.hpa", "HorizontalPodAutoscaler", "orders-api",
			corev1.EventTypeNormal, "SuccessfulRescale", "New size: 4; reason: cpu resource utilization above target", 90*time.Minute, 25*time.Minute, 3),
		event("production", "postgres-0.backup", "Pod", "postgres-0",
			corev1.EventTypeNormal, "SuccessfulCreate", "Created pod: nightly-backup-284910", 5*time.Hour, 5*time.Hour, 1),
		event("production", "postgres-1.disk", "Pod", "postgres-1",
			corev1.EventTypeWarning, "FailedMount", "MountVolume.SetUp failed for volume pvc: timeout waiting for volume (simulated)", 3*time.Hour, 2*time.Hour+40*time.Minute, 2),
		event("production", "checkout-oom", "Pod", "checkout-5b3c2-ee555",
			corev1.EventTypeWarning, "OOMKilling", "Memory cgroup out of memory (showcase warning)", 70*time.Minute, 70*time.Minute, 1),
		event("staging", "migration.completed", "Job", "db-migrate",
			corev1.EventTypeNormal, "Completed", "Job completed successfully", 26*time.Hour, 26*time.Hour, 1),
		event("staging", "web-frontend-image", "Pod", "web-frontend-9a1b2-gg777",
			corev1.EventTypeWarning, "Failed", "Failed to pull image: rate limited by registry (simulated)", 4*time.Hour, 3*time.Hour, 5),
		event("monitoring", "prometheus.rule", "Pod", "prometheus-0abc1",
			corev1.EventTypeNormal, "Saw", "Loaded 42 recording rules (showcase)", 12*time.Hour, time.Hour, 8),
		event("monitoring", "node-exporter-bbb.restart", "Pod", "node-exporter-bbb",
			corev1.EventTypeWarning, "BackOff", "Back-off restarting failed container exporter", 55*time.Minute, 20*time.Minute, 6),
		event("kube-system", "coredns-xyz01.sync", "Pod", "coredns-xyz01",
			corev1.EventTypeNormal, "Synced", "CoreDNS configmap synced", 48*time.Hour, 15*time.Minute, 12),
		event("cilibase", "cilikube-demo.ready", "Deployment", "cilikube-demo",
			corev1.EventTypeNormal, "Available", "Deployment has minimum availability", 30*time.Hour, 30*time.Hour, 1),
		event("default", "ingress.web", "Ingress", "web",
			corev1.EventTypeNormal, "Sync", "Scheduled for sync (showcase ingress controller)", 9*time.Hour, 40*time.Minute, 3),
		event("production", "deny-external.applied", "NetworkPolicy", "deny-external",
			corev1.EventTypeNormal, "PolicyApplied", "NetworkPolicy selected 4 pods", 14*time.Hour, 14*time.Hour, 1),
		event("default", "web-frontend-7d9f8b-def34.kill", "Pod", "web-frontend-7d9f8b-def34",
			corev1.EventTypeWarning, "Killing", "Stopping container web-frontend (rolling update)", 35*time.Minute, 35*time.Minute, 1),
		event("default", "web-frontend-7d9f8b-ghi56.ready", "Pod", "web-frontend-7d9f8b-ghi56",
			corev1.EventTypeNormal, "Ready", "Readiness probe succeeded", 30*time.Minute, 2*time.Minute, 7),
	}

	out := make([]runtime.Object, 0, len(evs))
	for i, e := range evs {
		// Ensure unique names for fake clientset
		if e.Name == "" {
			e.Name = fmt.Sprintf("showcase-event-%d", i)
		}
		out = append(out, e)
	}
	return out
}

// stampCreationTimes sets CreationTimestamp on typed metav1 objects that lack one.
func stampCreationTimes(objs []runtime.Object) {
	for i, obj := range objs {
		mo, ok := obj.(metav1.Object)
		if !ok {
			continue
		}
		if !mo.GetCreationTimestamp().Time.IsZero() {
			continue
		}
		// Spread ages: ~30m … ~10d based on index/name for a lived-in look.
		hours := 1 + (i*3+len(mo.GetName()))%240
		mo.SetCreationTimestamp(ago(time.Duration(hours) * time.Hour))
		if mo.GetUID() == "" {
			// leave empty — fake client assigns on create; list still works
		}
	}
}
