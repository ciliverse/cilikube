package store

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Labels 是一个 map[string]string 的自定义类型，用于 GORM 存储
type Labels map[string]string

// Value - 实现 driver.Valuer 接口，GORM 在写入时会调用
func (l Labels) Value() (driver.Value, error) {
	if l == nil {
		return nil, nil
	}
	return json.Marshal(l)
}

// Scan - 实现 sql.Scanner 接口，GORM 在读取时会调用
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

// Cluster 是数据库中 cluster 表的 GORM 模型，为企业级管理而设计
type Cluster struct {
	// --- 核心标识 ---
	// ID 是集群的唯一、不可变标识符，使用 UUID。这是系统内部的主键
	ID string `gorm:"type:varchar(36);primary_key" json:"id"`
	// Name 是用户定义的、易于记忆的集群显示名称，必须唯一
	Name string `gorm:"type:varchar(100);unique;not null" json:"name"`

	// --- 连接信息 ---
	// KubeconfigData 存储加密后的 kubeconfig 内容本身，而不是路径
	// 这使得应用完全独立于环境，具备极佳的可移植性。
	KubeconfigData []byte `gorm:"type:blob;not null" json:"-"`

	// --- 元数据与描述 ---
	// Description 是对集群用途、位置等的详细描述
	Description string `gorm:"type:text" json:"description"`
	// Provider 是云服务商或环境，例如 "aws", "gcp", "minikube", "on-premise"
	Provider string `gorm:"type:varchar(50)" json:"provider"`
	// Environment 标记集群的所属环境，如 "production", "staging", "development"
	Environment string `gorm:"type:varchar(50);index" json:"environment"`
	// Region 是集群所在的地理区域，如 "us-east-1", "asia-northeast1"
	Region string `gorm:"type:varchar(50)" json:"region"`
	// Version 存储探测到的 Kubernetes Master 版本号
	Version string `gorm:"type:varchar(20)" json:"version"`

	// --- 状态与标签 ---
	// Status 是管理员设置的集群状态，如 "Active", "Maintenance", "Inactive"
	Status string `gorm:"type:varchar(50);default:'Active'" json:"status"`
	// Labels 提供了灵活的键值对标签，用于分组、筛选和策略应用，这是企业级管理的关键特性
	Labels Labels `gorm:"type:json" json:"labels"`

	// --- 审计信息 ---
	// GORM 会自动管理 CreatedAt 和 UpdatedAt 时间戳
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
