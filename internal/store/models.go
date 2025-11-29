package store

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Labels is a custom type of map[string]string for GORM storage
type Labels map[string]string

// Value - implements driver.Valuer interface, called by GORM when writing
func (l Labels) Value() (driver.Value, error) {
	if l == nil {
		return nil, nil
	}
	return json.Marshal(l)
}

// Scan - implements sql.Scanner interface, called by GORM when reading
func (l *Labels) Scan(value interface{}) error {
	if value == nil {
		*l = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &l)
}

// Cluster is the GORM model for the cluster table in the database, designed for enterprise-level management
type Cluster struct {
	// --- Core Identity ---
	// ID is the unique, immutable identifier for the cluster, using UUID. This is the internal primary key
	ID string `gorm:"type:varchar(36);primary_key" json:"id"`
	// Name is the user-defined, memorable display name for the cluster, must be unique
	Name string `gorm:"type:varchar(100);unique;not null" json:"name"`

	// --- Connection Information ---
	// KubeconfigData stores the encrypted kubeconfig content itself, not the path
	// This makes the application completely environment-independent with excellent portability
	KubeconfigData []byte `gorm:"type:blob;not null" json:"-"`

	// --- Metadata and Description ---
	// Description is a detailed description of the cluster's purpose, location, etc.
	Description string `gorm:"type:text" json:"description"`
	// Provider is the cloud service provider or environment, e.g., "aws", "gcp", "minikube", "on-premise"
	Provider string `gorm:"type:varchar(50)" json:"provider"`
	// Environment marks the cluster's environment, such as "production", "staging", "development"
	Environment string `gorm:"type:varchar(50);index" json:"environment"`
	// Region is the geographical region where the cluster is located, e.g., "us-east-1", "asia-northeast1"
	Region string `gorm:"type:varchar(50)" json:"region"`
	// Version stores the detected Kubernetes Master version number
	Version string `gorm:"type:varchar(20)" json:"version"`

	// --- Status and Labels ---
	// Status is the cluster status set by administrators, such as "Active", "Maintenance", "Inactive"
	Status string `gorm:"type:varchar(50);default:'Active'" json:"status"`
	// Labels provides flexible key-value pair labels for grouping, filtering, and policy application, a key feature for enterprise management
	Labels Labels `gorm:"type:json" json:"labels"`

	// --- Audit Information ---
	// GORM automatically manages CreatedAt and UpdatedAt timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
