package models

import "time"

type CreateClusterRequest struct {
	Name           string `json:"name" binding:"required"`
	KubeconfigData string `json:"kubeconfigData" binding:"required"`
	Provider       string `json:"provider"`
	Description    string `json:"description"`
	Environment    string `json:"environment"`
	Region         string `json:"region"`
}

type UpdateClusterRequest struct {
	Name           string            `json:"name"`
	Provider       string            `json:"provider"`
	Description    string            `json:"description"`
	Environment    string            `json:"environment"`
	Region         string            `json:"region"`
	Status         string            `json:"status"`
	Labels         map[string]string `json:"labels"`
	KubeconfigData string            `json:"kubeconfigData,omitempty"`
}

type ClusterResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Provider    string            `json:"provider"`
	Description string            `json:"description"`
	Environment string            `json:"environment"`
	Region      string            `json:"region"`
	Version     string            `json:"version"`
	Status      string            `json:"status"`
	Source      string            `json:"source"`
	Labels      map[string]string `json:"labels"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type ClusterListResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Server      string `json:"server"`
	Version     string `json:"version"`
	Status      string `json:"status"`
	Source      string `json:"source"`
	Environment string `json:"environment"`
}
