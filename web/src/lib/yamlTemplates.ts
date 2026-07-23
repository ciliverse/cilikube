/** Minimal YAML skeletons for create-from-YAML flows. */
export function yamlTemplate(resource: string, namespace?: string): string {
  const ns = namespace || 'default'
  const stubs: Record<string, string> = {
    pods: `apiVersion: v1
kind: Pod
metadata:
  name: example-pod
  namespace: ${ns}
spec:
  containers:
    - name: main
      image: nginx:alpine
`,
    deployments: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: example-deploy
  namespace: ${ns}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: example
  template:
    metadata:
      labels:
        app: example
    spec:
      containers:
        - name: main
          image: nginx:alpine
`,
    statefulsets: `apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: example-sts
  namespace: ${ns}
spec:
  serviceName: example-sts
  replicas: 1
  selector:
    matchLabels:
      app: example-sts
  template:
    metadata:
      labels:
        app: example-sts
    spec:
      containers:
        - name: main
          image: nginx:alpine
`,
    daemonsets: `apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: example-ds
  namespace: ${ns}
spec:
  selector:
    matchLabels:
      app: example-ds
  template:
    metadata:
      labels:
        app: example-ds
    spec:
      containers:
        - name: main
          image: nginx:alpine
`,
    jobs: `apiVersion: batch/v1
kind: Job
metadata:
  name: example-job
  namespace: ${ns}
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: main
          image: busybox:1.36
          command: ["echo", "hello"]
`,
    cronjobs: `apiVersion: batch/v1
kind: CronJob
metadata:
  name: example-cron
  namespace: ${ns}
spec:
  schedule: "0 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: main
              image: busybox:1.36
              command: ["echo", "hello"]
`,
    services: `apiVersion: v1
kind: Service
metadata:
  name: example-svc
  namespace: ${ns}
spec:
  selector:
    app: example
  ports:
    - port: 80
      targetPort: 80
`,
    ingresses: `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example-ingress
  namespace: ${ns}
spec:
  rules:
    - host: example.local
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: example-svc
                port:
                  number: 80
`,
    configmaps: `apiVersion: v1
kind: ConfigMap
metadata:
  name: example-cm
  namespace: ${ns}
data:
  key: value
`,
    secrets: `apiVersion: v1
kind: Secret
metadata:
  name: example-secret
  namespace: ${ns}
type: Opaque
stringData:
  key: value
`,
    persistentvolumeclaims: `apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: example-pvc
  namespace: ${ns}
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
`,
    networkpolicies: `apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: example-netpol
  namespace: ${ns}
spec:
  podSelector: {}
  policyTypes:
    - Ingress
`,
    serviceaccounts: `apiVersion: v1
kind: ServiceAccount
metadata:
  name: example-sa
  namespace: ${ns}
`,
    roles: `apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: example-role
  namespace: ${ns}
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "list"]
`,
    rolebindings: `apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: example-rb
  namespace: ${ns}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: example-role
subjects:
  - kind: ServiceAccount
    name: default
    namespace: ${ns}
`,
    namespaces: `apiVersion: v1
kind: Namespace
metadata:
  name: example-ns
`,
    persistentvolumes: `apiVersion: v1
kind: PersistentVolume
metadata:
  name: example-pv
spec:
  capacity:
    storage: 1Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: /tmp/example-pv
`,
    storageclasses: `apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: example-sc
provisioner: kubernetes.io/hostpath
`,
    clusterroles: `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: example-clusterrole
rules:
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list"]
`,
    clusterrolebindings: `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: example-crb
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: example-clusterrole
subjects:
  - kind: ServiceAccount
    name: default
    namespace: default
`,
    horizontalpodautoscalers: `apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: example-hpa
  namespace: ${ns}
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: example-deploy
  minReplicas: 1
  maxReplicas: 3
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 80
`,
    poddisruptionbudgets: `apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: example-pdb
  namespace: ${ns}
spec:
  minAvailable: 1
  selector:
    matchLabels:
      app: example
`,
    resourcequotas: `apiVersion: v1
kind: ResourceQuota
metadata:
  name: example-quota
  namespace: ${ns}
spec:
  hard:
    requests.cpu: "4"
    requests.memory: 4Gi
    pods: "10"
`,
    limitranges: `apiVersion: v1
kind: LimitRange
metadata:
  name: example-limits
  namespace: ${ns}
spec:
  limits:
    - type: Container
      default:
        cpu: 200m
        memory: 256Mi
      defaultRequest:
        cpu: 50m
        memory: 64Mi
`,
  }

  return (
    stubs[resource] ||
    `apiVersion: v1
kind: ConfigMap
metadata:
  name: example
  namespace: ${ns}
data: {}
`
  )
}
