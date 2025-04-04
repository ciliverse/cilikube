package models

import (
	"encoding/base64" // Needed if handling binary data display/edit

	corev1 "k8s.io/api/core/v1"
)

// SecretResponse structures the data for a Secret sent to the UI list view.
// Does NOT include the actual secret data for security.
type SecretResponse struct {
	Name            string            `json:"name"`
	Namespace       string            `json:"namespace"`
	UID             string            `json:"uid"`
	Type            corev1.SecretType `json:"type"`      // e.g., Opaque, kubernetes.io/service-account-token
	DataCount       int               `json:"dataCount"` // Number of keys in the data/stringData field
	CreatedAt       string            `json:"createdAt"`
	Labels          map[string]string `json:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty"`
	ResourceVersion string            `json:"resourceVersion"`
}

// SecretListResponse structures the list response for Secrets.
type SecretListResponse struct {
	Items []SecretResponse `json:"items"`
	Total int              `json:"total"`
}

// ToSecretResponse converts a corev1.Secret to our API list response model.
func ToSecretResponse(secret *corev1.Secret) SecretResponse {
	dataCount := 0
	if secret.Data != nil {
		dataCount += len(secret.Data)
	}
	// StringData also contributes to the keys available
	if secret.StringData != nil {
		// Note: Keys in StringData might overlap with Data, but K8s handles the merge.
		// For a simple count, just add lengths. For precise unique keys, more logic needed.
		dataCount += len(secret.StringData)
	}

	return SecretResponse{
		Name:            secret.Name,
		Namespace:       secret.Namespace,
		UID:             string(secret.UID),
		Type:            secret.Type,
		DataCount:       dataCount,
		CreatedAt:       secret.CreationTimestamp.Format("2006-01-02T15:04:05Z"),
		Labels:          secret.Labels,
		Annotations:     secret.Annotations,
		ResourceVersion: secret.ResourceVersion,
	}
}

// SecretDetailResponse includes the actual data (use for Get/Edit via YAML)
// Data values are typically base64 encoded by K8s in the 'Data' field.
// We might decode StringData for easier editing in YAML view? Or keep raw.
type SecretDetailResponse struct {
	SecretResponse                   // Embed basic info
	Data           map[string]string `json:"data,omitempty"`       // Values will be base64 encoded strings
	StringData     map[string]string `json:"stringData,omitempty"` // Values are plain strings
	Type           corev1.SecretType `json:"type"`                 // Re-declare type from corev1 if not embedded fully
}

// ToSecretDetailResponse converts corev1.Secret to detailed response.
// It keeps Data base64 encoded as it comes from K8s API.
func ToSecretDetailResponse(secret *corev1.Secret) SecretDetailResponse {
	basicResponse := ToSecretResponse(secret)

	// Convert []byte from secret.Data to base64 strings for JSON response
	encodedData := make(map[string]string)
	if secret.Data != nil {
		for key, value := range secret.Data {
			encodedData[key] = base64.StdEncoding.EncodeToString(value)
		}
	}

	return SecretDetailResponse{
		SecretResponse: basicResponse,     // Basic info remains the same
		Data:           encodedData,       // Base64 encoded data
		StringData:     secret.StringData, // Plain string data
		Type:           secret.Type,       // Ensure type is present
	}
}
