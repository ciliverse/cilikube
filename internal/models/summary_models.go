package models

// ResourceSummary represents the count of various cluster resources.
// Use pointers to distinguish between a count of 0 and a failure to retrieve count.
type ResourceSummary struct {
	Nodes             *int `json:"nodes"`
	Namespaces        *int `json:"namespaces"`
	Pods              *int `json:"pods"`
	Deployments       *int `json:"deployments"`
	Services          *int `json:"services"`
	PersistentVolumes *int `json:"persistentVolumes"`
	Pvcs              *int `json:"pvcs"` // PersistentVolumeClaims
	StatefulSets      *int `json:"statefulSets"`
	DaemonSets        *int `json:"daemonSets"`
	ConfigMaps        *int `json:"configMaps"`
	Secrets           *int `json:"secrets"`
	Ingresses         *int `json:"ingresses"`
	// Add more resource types as needed
}
