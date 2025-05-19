package utils

import (
	"bytes"
	"fmt"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"regexp"
	"strconv"
)

var (
	dns1123Regex = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
)

// ValidateNamespace 验证命名空间格式
func ValidateNamespace(ns string) bool {
	return dns1123Regex.MatchString(ns) && len(ns) <= 63
}

// ValidateResourceName 验证资源名称格式
func ValidateResourceName(name string) bool {
	return dns1123Regex.MatchString(name) && len(name) <= 253
}

// ParseInt 安全转换字符串为整数
func ParseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}

// ParseDeploymentFromFile 解析 YAML/JSON 文件为 Deployment 对象（使用 Kubernetes 原生解码器）
func ParseDeploymentFromFile(data []byte) (*appsv1.Deployment, error) {
	// 使用 Kubernetes 提供的 YAML/JSON 解码器
	decoder := yaml.NewYAMLOrJSONDecoder(
		io.NopCloser(
			io.NewSectionReader(
				bytes.NewReader(data),
				0,
				int64(len(data)),
			),
		),
		1024,
	)

	var deployment appsv1.Deployment
	if err := decoder.Decode(&deployment); err != nil {
		return nil, fmt.Errorf("failed to decode YAML/JSON: %v", err.Error())
	}

	return &deployment, nil
}
