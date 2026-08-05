package service

import (
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

// Timeline health statuses (Radar-aligned).
const (
	TLHealthy   = "healthy"
	TLRolling   = "rolling"
	TLDegraded  = "degraded"
	TLUnhealthy = "unhealthy"
	TLUnknown   = "unknown"
)

// TimelineAppGroup extracts the app grouping label.
func TimelineAppGroup(labels map[string]string) string {
	if labels == nil {
		return "_ungrouped"
	}
	for _, k := range []string{"app.kubernetes.io/name", "app", "app.kubernetes.io/instance", "k8s-app"} {
		if v := strings.TrimSpace(labels[k]); v != "" {
			return v
		}
	}
	return "_ungrouped"
}

// ClassifyDeploymentStatus maps a Deployment to timeline status + reason.
func ClassifyDeploymentStatus(d *appsv1.Deployment) (string, string) {
	want := int32(1)
	if d.Spec.Replicas != nil {
		want = *d.Spec.Replicas
	}
	ready := d.Status.ReadyReplicas
	reason := fmt.Sprintf("%d/%d ready", ready, want)
	for _, c := range d.Status.Conditions {
		if c.Type == appsv1.DeploymentProgressing && c.Status == corev1.ConditionTrue &&
			strings.Contains(strings.ToLower(c.Reason), "progress") {
			if d.Status.UpdatedReplicas < want || d.Status.UnavailableReplicas > 0 {
				return TLRolling, reason
			}
		}
	}
	if d.Status.UnavailableReplicas > 0 {
		return TLUnhealthy, reason
	}
	if ready < want {
		if d.Status.UpdatedReplicas > 0 && d.Status.UpdatedReplicas < want {
			return TLRolling, reason
		}
		return TLDegraded, reason
	}
	return TLHealthy, reason
}

// ClassifyStatefulSetStatus maps a StatefulSet.
func ClassifyStatefulSetStatus(d *appsv1.StatefulSet) (string, string) {
	want := int32(1)
	if d.Spec.Replicas != nil {
		want = *d.Spec.Replicas
	}
	ready := d.Status.ReadyReplicas
	reason := fmt.Sprintf("%d/%d ready", ready, want)
	if d.Status.UpdatedReplicas > 0 && d.Status.UpdatedReplicas < want {
		return TLRolling, reason
	}
	if ready < want {
		return TLDegraded, reason
	}
	return TLHealthy, reason
}

// ClassifyDaemonSetStatus maps a DaemonSet.
func ClassifyDaemonSetStatus(d *appsv1.DaemonSet) (string, string) {
	want := d.Status.DesiredNumberScheduled
	ready := d.Status.NumberReady
	reason := fmt.Sprintf("%d/%d ready", ready, want)
	if d.Status.UpdatedNumberScheduled > 0 && d.Status.UpdatedNumberScheduled < want {
		return TLRolling, reason
	}
	if ready < want {
		return TLDegraded, reason
	}
	return TLHealthy, reason
}

// ClassifyPodStatus maps a Pod.
func ClassifyPodStatus(p *corev1.Pod) (string, string) {
	ready, total := 0, 0
	waiting := ""
	for _, cs := range p.Status.ContainerStatuses {
		total++
		if cs.Ready {
			ready++
		}
		if cs.State.Waiting != nil && cs.State.Waiting.Reason != "" {
			waiting = cs.State.Waiting.Reason
		}
	}
	reason := fmt.Sprintf("%d/%d %s", ready, total, p.Status.Phase)
	if waiting != "" {
		reason = fmt.Sprintf("%d/%d %s", ready, total, waiting)
		low := strings.ToLower(waiting)
		if strings.Contains(low, "crash") || strings.Contains(low, "backoff") ||
			strings.Contains(low, "err") || strings.Contains(low, "image") {
			return TLUnhealthy, reason
		}
	}
	switch p.Status.Phase {
	case corev1.PodFailed, corev1.PodUnknown:
		return TLUnhealthy, reason
	case corev1.PodPending:
		return TLDegraded, reason
	case corev1.PodSucceeded:
		return TLHealthy, reason
	}
	if total > 0 && ready < total {
		return TLDegraded, reason
	}
	if p.Status.Phase == corev1.PodRunning {
		return TLHealthy, reason
	}
	return TLUnknown, reason
}

// ClassifyJobStatus maps a Job.
func ClassifyJobStatus(j *batchv1.Job) (string, string) {
	reason := fmt.Sprintf("active=%d succeeded=%d failed=%d", j.Status.Active, j.Status.Succeeded, j.Status.Failed)
	for _, c := range j.Status.Conditions {
		if c.Type == batchv1.JobFailed && c.Status == corev1.ConditionTrue {
			return TLUnhealthy, reason
		}
		if c.Type == batchv1.JobComplete && c.Status == corev1.ConditionTrue {
			return TLHealthy, reason
		}
	}
	if j.Status.Active > 0 {
		return TLRolling, reason
	}
	if j.Status.Failed > 0 {
		return TLUnhealthy, reason
	}
	if j.Status.Succeeded > 0 {
		return TLHealthy, reason
	}
	return TLUnknown, reason
}

// ClassifyCronJobStatus maps a CronJob.
func ClassifyCronJobStatus(c *batchv1.CronJob) (string, string) {
	if c.Spec.Suspend != nil && *c.Spec.Suspend {
		return TLDegraded, "suspended"
	}
	return TLHealthy, c.Spec.Schedule
}

// ClassifyServiceStatus maps a Service (presence-only).
func ClassifyServiceStatus(svc *corev1.Service) (string, string) {
	if svc.Spec.ClusterIP == "" || svc.Spec.ClusterIP == "None" {
		return TLHealthy, string(svc.Spec.Type)
	}
	return TLHealthy, svc.Spec.ClusterIP
}

// TimelineHref builds console detail path.
func TimelineHref(kind, namespace, name string) string {
	k := strings.ToLower(kind)
	switch k {
	case "deployment":
		return fmt.Sprintf("/deployments/%s/%s", namespace, name)
	case "statefulset":
		return fmt.Sprintf("/statefulsets/%s/%s", namespace, name)
	case "daemonset":
		return fmt.Sprintf("/daemonsets/%s/%s", namespace, name)
	case "pod":
		return fmt.Sprintf("/pods/%s/%s", namespace, name)
	case "job":
		return fmt.Sprintf("/jobs/%s/%s", namespace, name)
	case "cronjob":
		return fmt.Sprintf("/cronjobs/%s/%s", namespace, name)
	case "service":
		return fmt.Sprintf("/services/%s/%s", namespace, name)
	default:
		return "/"
	}
}

// MapEventMarker converts a K8s event into a timeline marker kind.
func MapEventMarker(eventType, reason string) string {
	r := strings.ToLower(reason)
	if strings.Contains(r, "delet") || strings.Contains(r, "kill") {
		return "deleted"
	}
	if strings.EqualFold(eventType, "Warning") {
		return "warning"
	}
	if strings.Contains(r, "creat") || strings.Contains(r, "schedul") || r == "started" || r == "pulled" {
		return "created"
	}
	return "modified"
}
