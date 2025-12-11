package models

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CRDListResponse represents the response for CRD list
type CRDListResponse struct {
	Items []CRDItem `json:"items"`
	Total int       `json:"total"`
}

// CRDItem represents CRD item information
type CRDItem struct {
	Name        string            `json:"name"`
	Group       string            `json:"group"`
	Version     string            `json:"version"`
	Kind        string            `json:"kind"`
	Plural      string            `json:"plural"`
	Singular    string            `json:"singular"`
	Scope       string            `json:"scope"`
	Categories  []string          `json:"categories,omitempty"`
	ShortNames  []string          `json:"shortNames,omitempty"`
	CreatedAt   metav1.Time       `json:"createdAt"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// CustomResourceListResponse represents the response for custom resource list
type CustomResourceListResponse struct {
	Items []CustomResourceItem `json:"items"`
	Total int                  `json:"total"`
}

// CustomResourceItem represents a custom resource item
type CustomResourceItem struct {
	Name        string                 `json:"name"`
	Namespace   string                 `json:"namespace,omitempty"`
	Kind        string                 `json:"kind"`
	APIVersion  string                 `json:"apiVersion"`
	CreatedAt   metav1.Time            `json:"createdAt"`
	Labels      map[string]string      `json:"labels,omitempty"`
	Annotations map[string]string      `json:"annotations,omitempty"`
	Spec        map[string]interface{} `json:"spec,omitempty"`
	Status      map[string]interface{} `json:"status,omitempty"`
}

// CRDDetailResponse represents the response for CRD details
type CRDDetailResponse struct {
	CRDItem
	Schema      map[string]interface{} `json:"schema,omitempty"`
	Versions    []CRDVersion           `json:"versions,omitempty"`
	Conditions  []CRDCondition         `json:"conditions,omitempty"`
	Description string                 `json:"description,omitempty"`
}

// CRDVersion represents CRD version information
type CRDVersion struct {
	Name    string `json:"name"`
	Served  bool   `json:"served"`
	Storage bool   `json:"storage"`
}

// CRDCondition represents CRD status condition
type CRDCondition struct {
	Type               string      `json:"type"`
	Status             string      `json:"status"`
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	Reason             string      `json:"reason,omitempty"`
	Message            string      `json:"message,omitempty"`
}

// CustomResourceRequest represents the request for creating/updating custom resources
type CustomResourceRequest struct {
	APIVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   map[string]interface{} `json:"metadata"`
	Spec       map[string]interface{} `json:"spec,omitempty"`
}
