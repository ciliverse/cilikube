package store

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type gormClusterStore struct {
	db            *gorm.DB
	encryptionKey []byte
}

// NewGormClusterStore 创建一个新的 ClusterStore GORM 实现。
func NewGormClusterStore(db *gorm.DB, encryptionKey []byte) (ClusterStore, error) {
	if len(encryptionKey) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes long for AES-256")
	}
	return &gormClusterStore{
		db:            db,
		encryptionKey: encryptionKey,
	}, nil
}

func (s *gormClusterStore) CreateCluster(cluster *Cluster) error {
	cluster.ID = uuid.NewString()
	encryptedData, err := Encrypt(cluster.KubeconfigData, s.encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt kubeconfig: %w", err)
	}
	cluster.KubeconfigData = encryptedData
	return s.db.Create(cluster).Error
}

func (s *gormClusterStore) GetClusterByID(id string) (*Cluster, error) {
	var cluster Cluster
	if err := s.db.First(&cluster, "id = ?", id).Error; err != nil {
		return nil, err
	}
	decryptedData, err := Decrypt(cluster.KubeconfigData, s.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt kubeconfig for cluster %s: %w", cluster.Name, err)
	}
	cluster.KubeconfigData = decryptedData
	return &cluster, nil
}

func (s *gormClusterStore) GetClusterByName(name string) (*Cluster, error) {
	var cluster Cluster
	if err := s.db.First(&cluster, "name = ?", name).Error; err != nil {
		return nil, err
	}
	decryptedData, err := Decrypt(cluster.KubeconfigData, s.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt kubeconfig for cluster %s: %w", cluster.Name, err)
	}
	cluster.KubeconfigData = decryptedData
	return &cluster, nil
}

func (s *gormClusterStore) GetAllClusters() ([]Cluster, error) {
	var clusters []Cluster
	if err := s.db.Find(&clusters).Error; err != nil {
		return nil, err
	}
	for i := range clusters {
		decryptedData, err := Decrypt(clusters[i].KubeconfigData, s.encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt kubeconfig for cluster %s: %w", clusters[i].Name, err)
		}
		clusters[i].KubeconfigData = decryptedData
	}
	return clusters, nil
}

func (s *gormClusterStore) UpdateCluster(cluster *Cluster) error {
	if len(cluster.KubeconfigData) > 0 {
		encryptedData, err := Encrypt(cluster.KubeconfigData, s.encryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt kubeconfig for update on cluster %s: %w", cluster.Name, err)
		}
		cluster.KubeconfigData = encryptedData
	}
	return s.db.Save(cluster).Error
}

func (s *gormClusterStore) DeleteClusterByName(name string) error {
	return s.db.Delete(&Cluster{}, "name = ?", name).Error
}
