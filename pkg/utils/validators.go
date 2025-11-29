package utils

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var (
	dns1123Regex = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
)

// ValidateNamespace validates namespace format
func ValidateNamespace(ns string) bool {
	return dns1123Regex.MatchString(ns) && len(ns) <= 63
}

// ValidateResourceName validates resource name format
func ValidateResourceName(name string) bool {
	return dns1123Regex.MatchString(name) && len(name) <= 253
}

// ParseInt safely converts string to integer
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

// ParseDeploymentFromFile parses YAML/JSON file to Deployment object (using Kubernetes native decoder)
func ParseDeploymentFromFile(data []byte) (*appsv1.Deployment, error) {
	// Use YAML/JSON decoder provided by Kubernetes
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
