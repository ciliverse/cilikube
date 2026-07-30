package k8s

import (
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func nodeReady(name, role, cpu, mem string, ready bool) *corev1.Node {
	n := node(name, role, cpu, mem)
	status := corev1.ConditionTrue
	reason := "KubeletReady"
	msg := "showcase ready"
	if !ready {
		status = corev1.ConditionFalse
		reason = "KubeletNotReady"
		msg = "showcase simulated NotReady (disk pressure)"
	}
	n.Status.Conditions = []corev1.NodeCondition{
		{Type: corev1.NodeReady, Status: status, Reason: reason, Message: msg},
	}
	return n
}

func crashLoopPod(nsName, name, app, nodeName string) *corev1.Pod {
	p := pod(nsName, name, app, nodeName, string(corev1.PodRunning))
	p.Status.Phase = corev1.PodRunning
	p.Status.ContainerStatuses = []corev1.ContainerStatus{{
		Name:         app,
		Ready:        false,
		RestartCount: 12,
		Image:        "ghcr.io/ciliverse/" + app + ":showcase",
		State: corev1.ContainerState{
			Waiting: &corev1.ContainerStateWaiting{
				Reason:  "CrashLoopBackOff",
				Message: "back-off restarting failed container (showcase)",
			},
		},
	}}
	return p
}

// ShowcaseProdSeedObjects simulates a busier production-east cluster with some faults.
func ShowcaseProdSeedObjects() []runtime.Object {
	objs := []runtime.Object{
		ns("default"),
		ns("kube-system"),
		ns("production"),
		ns("payments"),
		ns("observability"),

		nodeReady("prod-cp-1", "control-plane", "8", "16Gi", true),
		nodeReady("prod-worker-1", "worker", "16", "64Gi", true),
		nodeReady("prod-worker-2", "worker", "16", "64Gi", true),
		nodeReady("prod-worker-3", "worker", "16", "64Gi", false),

		deploy("production", "orders-api", 5),
		deploy("production", "checkout", 3),
		deploy("production", "inventory", 2),
		svc("production", "orders-api", 8080),
		svc("production", "checkout", 8081),
		svc("production", "inventory", 8082),

		pod("production", "orders-api-a1b2c-11111", "orders-api", "prod-worker-1", "Running"),
		pod("production", "orders-api-a1b2c-22222", "orders-api", "prod-worker-2", "Running"),
		pod("production", "orders-api-a1b2c-33333", "orders-api", "prod-worker-1", "Running"),
		pod("production", "orders-api-a1b2c-44444", "orders-api", "prod-worker-2", "Running"),
		crashLoopPod("production", "orders-api-a1b2c-55555", "orders-api", "prod-worker-3"),
		pod("production", "checkout-c9d8e-aaaa1", "checkout", "prod-worker-1", "Running"),
		pod("production", "checkout-c9d8e-bbbb2", "checkout", "prod-worker-2", "Running"),
		pod("production", "checkout-c9d8e-cccc3", "checkout", "prod-worker-1", "Pending"),
		pod("production", "inventory-f1e2d-dddd1", "inventory", "prod-worker-2", "Running"),
		pod("production", "inventory-f1e2d-eeee2", "inventory", "prod-worker-1", "Failed"),

		deploy("payments", "pay-gateway", 3),
		svc("payments", "pay-gateway", 8443),
		pod("payments", "pay-gateway-9x8y7-p1111", "pay-gateway", "prod-worker-1", "Running"),
		pod("payments", "pay-gateway-9x8y7-p2222", "pay-gateway", "prod-worker-2", "Running"),
		crashLoopPod("payments", "pay-gateway-9x8y7-p3333", "pay-gateway", "prod-worker-3"),

		deploy("observability", "prometheus", 1),
		pod("observability", "prometheus-0", "prometheus", "prod-worker-1", "Running"),
		svc("observability", "prometheus", 9090),

		deploy("kube-system", "coredns", 2),
		pod("kube-system", "coredns-prod01", "coredns", "prod-cp-1", "Running"),
		pod("kube-system", "coredns-prod02", "coredns", "prod-cp-1", "Running"),
		svc("kube-system", "kube-dns", 53),

		&appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{Name: "node-exporter", Namespace: "observability", Labels: appLabels("node-exporter")},
			Spec: appsv1.DaemonSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: appLabels("node-exporter")},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{Labels: appLabels("node-exporter")},
					Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "exporter", Image: "prom/node-exporter:showcase"}}},
				},
			},
			Status: appsv1.DaemonSetStatus{DesiredNumberScheduled: 4, NumberReady: 3, NumberAvailable: 3},
		},
		pod("observability", "node-exporter-cp", "node-exporter", "prod-cp-1", "Running"),
		pod("observability", "node-exporter-w1", "node-exporter", "prod-worker-1", "Running"),
		pod("observability", "node-exporter-w2", "node-exporter", "prod-worker-2", "Running"),
		pod("observability", "node-exporter-w3", "node-exporter", "prod-worker-3", "Pending"),
	}

	objs = append(objs, ShowcaseProdEventObjects()...)
	stampCreationTimes(objs)
	return objs
}

// ShowcaseStagingSeedObjects simulates a smaller staging lab (mostly healthy).
func ShowcaseStagingSeedObjects() []runtime.Object {
	objs := []runtime.Object{
		ns("default"),
		ns("kube-system"),
		ns("staging"),
		ns("preview"),

		nodeReady("stg-cp-1", "control-plane", "4", "8Gi", true),
		nodeReady("stg-worker-1", "worker", "8", "16Gi", true),

		deploy("staging", "web-frontend", 2),
		deploy("staging", "api-gateway", 1),
		svc("staging", "web-frontend", 80),
		svc("staging", "api-gateway", 8080),
		pod("staging", "web-frontend-stg-aaa01", "web-frontend", "stg-worker-1", "Running"),
		pod("staging", "web-frontend-stg-bbb02", "web-frontend", "stg-worker-1", "Running"),
		pod("staging", "api-gateway-stg-ccc03", "api-gateway", "stg-worker-1", "Running"),

		deploy("preview", "feature-x", 1),
		svc("preview", "feature-x", 8080),
		pod("preview", "feature-x-pr42-ddd04", "feature-x", "stg-worker-1", "Running"),
		pod("preview", "feature-x-pr42-eee05", "feature-x", "stg-worker-1", "Pending"),

		deploy("kube-system", "coredns", 1),
		pod("kube-system", "coredns-stg01", "coredns", "stg-cp-1", "Running"),
		svc("kube-system", "kube-dns", 53),
	}

	objs = append(objs, ShowcaseStagingEventObjects()...)
	stampCreationTimes(objs)
	return objs
}

// ShowcaseProdEventObjects — denser Warning stream for fleet / AI drill-down demos.
func ShowcaseProdEventObjects() []runtime.Object {
	evs := []*corev1.Event{
		event("production", "orders-api-55555.crash", "Pod", "orders-api-a1b2c-55555",
			corev1.EventTypeWarning, "BackOff", "Back-off restarting failed container orders-api", 40*time.Minute, 3*time.Minute, 14),
		event("production", "orders-api-55555.oom", "Pod", "orders-api-a1b2c-55555",
			corev1.EventTypeWarning, "OOMKilling", "Memory cgroup out of memory (showcase)", 45*time.Minute, 8*time.Minute, 3),
		event("production", "checkout-cccc3.schedule", "Pod", "checkout-c9d8e-cccc3",
			corev1.EventTypeWarning, "FailedScheduling", "0/4 nodes available: 1 node(s) NotReady (showcase)", 25*time.Minute, 5*time.Minute, 8),
		event("production", "inventory-eeee2.fail", "Pod", "inventory-f1e2d-eeee2",
			corev1.EventTypeWarning, "Failed", "Container exited with code 1 (showcase)", 2*time.Hour, 90*time.Minute, 2),
		event("payments", "pay-gateway-p3333.crash", "Pod", "pay-gateway-9x8y7-p3333",
			corev1.EventTypeWarning, "BackOff", "CrashLoopBackOff (showcase payment probe)", 55*time.Minute, 4*time.Minute, 11),
		event("payments", "pay-gateway-p3333.unhealthy", "Pod", "pay-gateway-9x8y7-p3333",
			corev1.EventTypeWarning, "Unhealthy", "Liveness probe failed: HTTP 500 (showcase)", 50*time.Minute, 6*time.Minute, 9),
		event("observability", "node-exporter-w3.pending", "Pod", "node-exporter-w3",
			corev1.EventTypeWarning, "FailedScheduling", "pod did not fit: node NotReady", 30*time.Minute, 10*time.Minute, 5),
		event("production", "orders-api.hpa", "HorizontalPodAutoscaler", "orders-api",
			corev1.EventTypeNormal, "SuccessfulRescale", "New size: 5; reason: cpu above target", 3*time.Hour, 40*time.Minute, 4),
		event("kube-system", "coredns-prod01.sync", "Pod", "coredns-prod01",
			corev1.EventTypeNormal, "Synced", "CoreDNS configmap synced", 48*time.Hour, 20*time.Minute, 9),
		event("production", "checkout-aaaa1.ready", "Pod", "checkout-c9d8e-aaaa1",
			corev1.EventTypeNormal, "Started", "Container checkout started", 6*time.Hour, 6*time.Hour, 1),
		event("payments", "pay-gateway-p1111.ready", "Pod", "pay-gateway-9x8y7-p1111",
			corev1.EventTypeNormal, "Pulled", "Successfully pulled image pay-gateway:showcase", 5*time.Hour, 5*time.Hour, 1),
		event("observability", "prometheus.rule", "Pod", "prometheus-0",
			corev1.EventTypeWarning, "FailedGetResourceMetric", "unable to get metric cpu: no metrics for showcase window", 70*time.Minute, 15*time.Minute, 6),
	}
	out := make([]runtime.Object, 0, len(evs))
	for _, e := range evs {
		out = append(out, e)
	}
	return out
}

// ShowcaseStagingEventObjects — light warning noise for staging.
func ShowcaseStagingEventObjects() []runtime.Object {
	evs := []*corev1.Event{
		event("staging", "web-frontend.scaled", "Deployment", "web-frontend",
			corev1.EventTypeNormal, "ScalingReplicaSet", "Scaled up replica set to 2", 8*time.Hour, 8*time.Hour, 1),
		event("preview", "feature-x-eee05.image", "Pod", "feature-x-pr42-eee05",
			corev1.EventTypeWarning, "Failed", "Failed to pull image: registry rate limited (showcase)", 90*time.Minute, 20*time.Minute, 4),
		event("staging", "api-gateway.ready", "Pod", "api-gateway-stg-ccc03",
			corev1.EventTypeNormal, "Started", "Container api-gateway started", 4*time.Hour, 4*time.Hour, 1),
		event("kube-system", "coredns-stg01.sync", "Pod", "coredns-stg01",
			corev1.EventTypeNormal, "Synced", "CoreDNS configmap synced", 24*time.Hour, 30*time.Minute, 3),
	}
	out := make([]runtime.Object, 0, len(evs))
	for _, e := range evs {
		out = append(out, e)
	}
	return out
}
