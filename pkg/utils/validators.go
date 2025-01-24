package utils

import (
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
