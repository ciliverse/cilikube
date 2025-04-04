package models

import corev1 "k8s.io/api/core/v1"

// ConfigMapResponse structures the data for a ConfigMap sent to the UI.
// In the list view, we might only show the number of keys.
type ConfigMapResponse struct {
	Name            string            `json:"name"`
	Namespace       string            `json:"namespace"`
	UID             string            `json:"uid"`
	DataCount       int               `json:"dataCount"` // Number of keys in the data field
	CreatedAt       string            `json:"createdAt"`
	Labels          map[string]string `json:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty"`
	ResourceVersion string            `json:"resourceVersion"`
	// Include full Data map only when getting a single item, not in the list response?
	// Data            map[string]string `json:"data,omitempty"`
}

// ConfigMapListResponse structures the list response for ConfigMaps.
type ConfigMapListResponse struct {
	Items []ConfigMapResponse `json:"items"`
	Total int                 `json:"total"`
}

// ToConfigMapResponse converts a corev1.ConfigMap to our API response model.
func ToConfigMapResponse(cm *corev1.ConfigMap) ConfigMapResponse {
	dataCount := 0
	if cm.Data != nil {
		dataCount = len(cm.Data)
	}
	// BinaryData count could also be added if needed:
	// if cm.BinaryData != nil { dataCount += len(cm.BinaryData) }

	return ConfigMapResponse{
		Name:            cm.Name,
		Namespace:       cm.Namespace,
		UID:             string(cm.UID),
		DataCount:       dataCount,
		CreatedAt:       cm.CreationTimestamp.Format("2006-01-02T15:04:05Z"),
		Labels:          cm.Labels,
		Annotations:     cm.Annotations,
		ResourceVersion: cm.ResourceVersion,
	}
}

// ToConfigMapDetailResponse includes the actual data (use for Get)
// You can merge this with ConfigMapResponse if always sending data
type ConfigMapDetailResponse struct {
	ConfigMapResponse                   // Embed basic info
	Data              map[string]string `json:"data,omitempty"`
	BinaryData        map[string][]byte `json:"binaryData,omitempty"` // Base64 encoded in JSON
}

func ToConfigMapDetailResponse(cm *corev1.ConfigMap) ConfigMapDetailResponse {
	basicResponse := ToConfigMapResponse(cm)
	return ConfigMapDetailResponse{
		ConfigMapResponse: basicResponse,
		Data:              cm.Data,
		BinaryData:        cm.BinaryData,
	}
}
