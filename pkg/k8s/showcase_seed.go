package k8s

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func i32(n int32) *int32 { return &n }

func appLabels(app string) map[string]string {
	return map[string]string{"app": app, "showcase": "true"}
}

func qty(s string) resource.Quantity {
	return resource.MustParse(s)
}

func ns(name string) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: name, Labels: map[string]string{"showcase": "true"}},
	}
}

func node(name, role string, cpu, mem string) *corev1.Node {
	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"kubernetes.io/hostname":              name,
				"node-role.kubernetes.io/" + role:     "",
				"topology.kubernetes.io/region":       "demo-east",
				"showcase":                            "true",
			},
		},
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{
				{Type: corev1.NodeReady, Status: corev1.ConditionTrue, Reason: "KubeletReady", Message: "showcase ready"},
			},
			Capacity: corev1.ResourceList{
				corev1.ResourceCPU:              qty(cpu),
				corev1.ResourceMemory:           qty(mem),
				corev1.ResourcePods:             qty("110"),
				corev1.ResourceEphemeralStorage: qty("100Gi"),
			},
			Allocatable: corev1.ResourceList{
				corev1.ResourceCPU:              qty(cpu),
				corev1.ResourceMemory:           qty(mem),
				corev1.ResourcePods:             qty("110"),
				corev1.ResourceEphemeralStorage: qty("90Gi"),
			},
			Addresses: []corev1.NodeAddress{
				{Type: corev1.NodeInternalIP, Address: fmt.Sprintf("10.42.0.%d", 10+len(name)%40)},
				{Type: corev1.NodeHostName, Address: name},
			},
			NodeInfo: corev1.NodeSystemInfo{
				KubeletVersion:          ShowcaseVersion,
				ContainerRuntimeVersion: "containerd://1.7.0-showcase",
				OperatingSystem:         "linux",
				Architecture:            "amd64",
			},
		},
	}
}

func pod(ns, name, app, nodeName, phase string, containers ...string) *corev1.Pod {
	if len(containers) == 0 {
		containers = []string{app}
	}
	cs := make([]corev1.Container, 0, len(containers))
	statuses := make([]corev1.ContainerStatus, 0, len(containers))
	for _, c := range containers {
		cs = append(cs, corev1.Container{
			Name:  c,
			Image: fmt.Sprintf("ghcr.io/ciliverse/%s:showcase", c),
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{corev1.ResourceCPU: qty("50m"), corev1.ResourceMemory: qty("64Mi")},
				Limits:   corev1.ResourceList{corev1.ResourceCPU: qty("500m"), corev1.ResourceMemory: qty("256Mi")},
			},
		})
		statuses = append(statuses, corev1.ContainerStatus{
			Name:  c,
			Ready: phase == string(corev1.PodRunning),
			State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{StartedAt: metav1.Now()}},
			Image: fmt.Sprintf("ghcr.io/ciliverse/%s:showcase", c),
		})
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels:    appLabels(app),
		},
		Spec: corev1.PodSpec{
			NodeName:   nodeName,
			Containers: cs,
		},
		Status: corev1.PodStatus{
			Phase:             corev1.PodPhase(phase),
			ContainerStatuses: statuses,
			PodIP:             "10.42.1." + fmt.Sprint((len(name)*7)%200+10),
			HostIP:            "10.42.0.11",
		},
	}
}

func deploy(ns, name string, replicas int32) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns, Labels: appLabels(name)},
		Spec: appsv1.DeploymentSpec{
			Replicas: i32(replicas),
			Selector: &metav1.LabelSelector{MatchLabels: appLabels(name)},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: appLabels(name)},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  name,
						Image: fmt.Sprintf("ghcr.io/ciliverse/%s:showcase", name),
					}},
				},
			},
		},
		Status: appsv1.DeploymentStatus{
			Replicas:          replicas,
			ReadyReplicas:     replicas,
			AvailableReplicas: replicas,
			UpdatedReplicas:   replicas,
		},
	}
}

func svc(ns, name string, port int32) *corev1.Service {
	return svcApp(ns, name, name, port)
}

// svcApp creates a Service whose selector matches appLabels(app) (e.g. kube-dns → coredns).
func svcApp(ns, name, app string, port int32) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns, Labels: appLabels(app)},
		Spec: corev1.ServiceSpec{
			Selector:  appLabels(app),
			Ports:     []corev1.ServicePort{{Port: port, TargetPort: intstr.FromInt32(port), Protocol: corev1.ProtocolTCP}},
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: fmt.Sprintf("10.96.%d.%d", port%200, (len(name)*3)%200),
		},
	}
}

// ShowcaseSeedObjects returns a rich fake inventory for every common resource kind.
func ShowcaseSeedObjects() []runtime.Object {
	objs := []runtime.Object{
		ns("default"),
		ns("kube-system"),
		ns("production"),
		ns("staging"),
		ns("monitoring"),
		ns("cilibase"),

		node("demo-master-1", "control-plane", "4", "8Gi"),
		node("demo-worker-1", "worker", "8", "16Gi"),
		node("demo-worker-2", "worker", "8", "16Gi"),

		// --- default ---
		deploy("default", "web-frontend", 3),
		deploy("default", "api-gateway", 2),
		svc("default", "web-frontend", 80),
		svc("default", "api-gateway", 8080),
		pod("default", "web-frontend-7d9f8b-abc12", "web-frontend", "demo-worker-1", "Running"),
		pod("default", "web-frontend-7d9f8b-def34", "web-frontend", "demo-worker-2", "Running"),
		pod("default", "web-frontend-7d9f8b-ghi56", "web-frontend", "demo-worker-1", "Running"),
		pod("default", "api-gateway-6c4d5-jkl78", "api-gateway", "demo-worker-2", "Running"),
		pod("default", "api-gateway-6c4d5-mno90", "api-gateway", "demo-worker-1", "Running"),

		&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "app-config", Namespace: "default", Labels: appLabels("web-frontend")},
			Data:       map[string]string{"APP_ENV": "showcase", "FEATURE_GLOBE": "on"},
		},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "app-secrets", Namespace: "default", Labels: appLabels("api-gateway")},
			Type:       corev1.SecretTypeOpaque,
			Data:       map[string][]byte{"token": []byte("showcase-not-a-real-secret")},
		},

		// --- production ---
		deploy("production", "orders-api", 4),
		deploy("production", "checkout", 2),
		svc("production", "orders-api", 8080),
		svc("production", "checkout", 8081),
		pod("production", "orders-api-8f2a1-aa111", "orders-api", "demo-worker-1", "Running"),
		pod("production", "orders-api-8f2a1-bb222", "orders-api", "demo-worker-2", "Running"),
		pod("production", "orders-api-8f2a1-cc333", "orders-api", "demo-worker-1", "Running"),
		pod("production", "orders-api-8f2a1-dd444", "orders-api", "demo-worker-2", "Running"),
		pod("production", "checkout-5b3c2-ee555", "checkout", "demo-worker-1", "Running"),
		pod("production", "checkout-5b3c2-ff666", "checkout", "demo-worker-2", "Running"),
		&appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{Name: "postgres", Namespace: "production", Labels: appLabels("postgres")},
			Spec: appsv1.StatefulSetSpec{
				Replicas: i32(2),
				Selector: &metav1.LabelSelector{MatchLabels: appLabels("postgres")},
				ServiceName: "postgres",
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{Labels: appLabels("postgres")},
					Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: "postgres", Image: "postgres:16-showcase"}}},
				},
			},
			Status: appsv1.StatefulSetStatus{Replicas: 2, ReadyReplicas: 2},
		},
		pod("production", "postgres-0", "postgres", "demo-worker-1", "Running"),
		pod("production", "postgres-1", "postgres", "demo-worker-2", "Running"),
		svc("production", "postgres", 5432),

		// --- staging ---
		deploy("staging", "web-frontend", 1),
		pod("staging", "web-frontend-9a1b2-gg777", "web-frontend", "demo-worker-2", "Running"),
		pod("staging", "migration-job-hh888", "migration", "demo-worker-1", "Succeeded"),
		svc("staging", "web-frontend", 80),

		// --- monitoring ---
		&appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{Name: "node-exporter", Namespace: "monitoring", Labels: appLabels("node-exporter")},
			Spec: appsv1.DaemonSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: appLabels("node-exporter")},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{Labels: appLabels("node-exporter")},
					Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: "exporter", Image: "prom/node-exporter:showcase"}}},
				},
			},
			Status: appsv1.DaemonSetStatus{DesiredNumberScheduled: 3, NumberReady: 3, NumberAvailable: 3},
		},
		pod("monitoring", "node-exporter-aaa", "node-exporter", "demo-master-1", "Running"),
		pod("monitoring", "node-exporter-bbb", "node-exporter", "demo-worker-1", "Running"),
		pod("monitoring", "node-exporter-ccc", "node-exporter", "demo-worker-2", "Running"),
		deploy("monitoring", "prometheus", 1),
		pod("monitoring", "prometheus-0abc1", "prometheus", "demo-worker-1", "Running"),
		svc("monitoring", "prometheus", 9090),

		// --- kube-system ---
		deploy("kube-system", "coredns", 2),
		pod("kube-system", "coredns-xyz01", "coredns", "demo-master-1", "Running"),
		pod("kube-system", "coredns-xyz02", "coredns", "demo-master-1", "Running"),
		svcApp("kube-system", "kube-dns", "coredns", 53),
		&appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{Name: "kube-proxy", Namespace: "kube-system", Labels: appLabels("kube-proxy")},
			Spec: appsv1.DaemonSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: appLabels("kube-proxy")},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{Labels: appLabels("kube-proxy")},
					Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: "kube-proxy", Image: "kube-proxy:showcase"}}},
				},
			},
			Status: appsv1.DaemonSetStatus{DesiredNumberScheduled: 3, NumberReady: 3, NumberAvailable: 3},
		},

		// Jobs / CronJobs
		&batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{Name: "db-migrate", Namespace: "staging", Labels: appLabels("migration")},
			Spec: batchv1.JobSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						RestartPolicy: corev1.RestartPolicyNever,
						Containers:    []corev1.Container{{Name: "migrate", Image: "migrate:showcase"}},
					},
				},
			},
			Status: batchv1.JobStatus{Succeeded: 1},
		},
		&batchv1.CronJob{
			ObjectMeta: metav1.ObjectMeta{Name: "nightly-backup", Namespace: "production", Labels: appLabels("backup")},
			Spec: batchv1.CronJobSpec{
				Schedule: "0 2 * * *",
				JobTemplate: batchv1.JobTemplateSpec{
					Spec: batchv1.JobSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								RestartPolicy: corev1.RestartPolicyOnFailure,
								Containers:    []corev1.Container{{Name: "backup", Image: "backup:showcase"}},
							},
						},
					},
				},
			},
		},

		// Storage
		&storagev1.StorageClass{
			ObjectMeta:  metav1.ObjectMeta{Name: "demo-ssd", Labels: map[string]string{"showcase": "true"}},
			Provisioner: "showcase.cilikube.io/ssd",
		},
		&corev1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{Name: "pv-postgres-0", Labels: map[string]string{"showcase": "true"}},
			Spec: corev1.PersistentVolumeSpec{
				Capacity:                      corev1.ResourceList{corev1.ResourceStorage: qty("20Gi")},
				AccessModes:                   []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimRetain,
				StorageClassName:              "demo-ssd",
			},
			Status: corev1.PersistentVolumeStatus{Phase: corev1.VolumeBound},
		},
		&corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{Name: "postgres-data-postgres-0", Namespace: "production", Labels: appLabels("postgres")},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				StorageClassName: strPtr("demo-ssd"),
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{corev1.ResourceStorage: qty("20Gi")},
				},
			},
			Status: corev1.PersistentVolumeClaimStatus{Phase: corev1.ClaimBound},
		},

		// Networking
		&networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{Name: "web", Namespace: "default", Labels: appLabels("web-frontend")},
			Spec: networkingv1.IngressSpec{
				Rules: []networkingv1.IngressRule{{
					Host: "demo.cilikube.local",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{{
								Path:     "/",
								PathType: pathTypePrefix(),
								Backend: networkingv1.IngressBackend{
									Service: &networkingv1.IngressServiceBackend{
										Name: "web-frontend",
										Port: networkingv1.ServiceBackendPort{Number: 80},
									},
								},
							}},
						},
					},
				}},
			},
		},
		&networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{Name: "deny-external", Namespace: "production", Labels: appLabels("orders-api")},
			Spec: networkingv1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{MatchLabels: appLabels("orders-api")},
				PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeIngress},
			},
		},

		// RBAC
		&corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: "app", Namespace: "production"}},
		&corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: "default", Namespace: "default"}},
		&rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{Name: "pod-reader", Namespace: "production"},
			Rules:      []rbacv1.PolicyRule{{APIGroups: []string{""}, Resources: []string{"pods"}, Verbs: []string{"get", "list"}}},
		},
		&rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{Name: "app-pod-reader", Namespace: "production"},
			Subjects:   []rbacv1.Subject{{Kind: "ServiceAccount", Name: "app", Namespace: "production"}},
			RoleRef:    rbacv1.RoleRef{APIGroup: "rbac.authorization.k8s.io", Kind: "Role", Name: "pod-reader"},
		},
		&rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{Name: "showcase-view"},
			Rules:      []rbacv1.PolicyRule{{APIGroups: []string{"*"}, Resources: []string{"*"}, Verbs: []string{"get", "list", "watch"}}},
		},
		&rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{Name: "showcase-view-binding"},
			Subjects:   []rbacv1.Subject{{Kind: "Group", Name: "system:authenticated"}},
			RoleRef:    rbacv1.RoleRef{APIGroup: "rbac.authorization.k8s.io", Kind: "ClusterRole", Name: "showcase-view"},
		},

		// HPA / PDB / Quota / LimitRange
		&autoscalingv2.HorizontalPodAutoscaler{
			ObjectMeta: metav1.ObjectMeta{Name: "orders-api", Namespace: "production", Labels: appLabels("orders-api")},
			Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{Kind: "Deployment", Name: "orders-api", APIVersion: "apps/v1"},
				MinReplicas:    i32(2),
				MaxReplicas:    10,
			},
			Status: autoscalingv2.HorizontalPodAutoscalerStatus{CurrentReplicas: 4, DesiredReplicas: 4},
		},
		&policyv1.PodDisruptionBudget{
			ObjectMeta: metav1.ObjectMeta{Name: "orders-api-pdb", Namespace: "production"},
			Spec: policyv1.PodDisruptionBudgetSpec{
				MinAvailable: &intstr.IntOrString{Type: intstr.Int, IntVal: 2},
				Selector:     &metav1.LabelSelector{MatchLabels: appLabels("orders-api")},
			},
		},
		&corev1.ResourceQuota{
			ObjectMeta: metav1.ObjectMeta{Name: "team-quota", Namespace: "production"},
			Spec: corev1.ResourceQuotaSpec{
				Hard: corev1.ResourceList{
					corev1.ResourceRequestsCPU:    qty("20"),
					corev1.ResourceRequestsMemory: qty("40Gi"),
					corev1.ResourcePods:           qty("50"),
				},
			},
		},
		&corev1.LimitRange{
			ObjectMeta: metav1.ObjectMeta{Name: "defaults", Namespace: "staging"},
			Spec: corev1.LimitRangeSpec{
				Limits: []corev1.LimitRangeItem{{
					Type: corev1.LimitTypeContainer,
					Default: corev1.ResourceList{
						corev1.ResourceCPU:    qty("200m"),
						corev1.ResourceMemory: qty("256Mi"),
					},
				}},
			},
		},

		// cilibase sample
		deploy("cilibase", "cilikube-demo", 2),
		pod("cilibase", "cilikube-demo-1a2b3", "cilikube-demo", "demo-worker-1", "Running"),
		pod("cilibase", "cilikube-demo-4c5d6", "cilikube-demo", "demo-worker-2", "Running"),
		svc("cilibase", "cilikube-demo", 8080),
	}

	objs = append(objs, ShowcaseEventObjects()...)
	stampCreationTimes(objs)
	return objs
}

func strPtr(s string) *string { return &s }

func pathTypePrefix() *networkingv1.PathType {
	p := networkingv1.PathTypePrefix
	return &p
}
