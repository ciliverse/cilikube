package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogicalPermissionsFromPoliciesFullAccess(t *testing.T) {
	policies := [][]string{
		{"admin", "/api/v1/*", "*"},
	}
	perms := logicalPermissionsFromPolicies(policies)
	assert.Equal(t, KnownLogicalPermissions(), perms)
}

func TestLogicalPermissionsFromPoliciesPartial(t *testing.T) {
	policies := [][]string{
		{"viewer", "/api/v1/pods/*", "GET"},
		{"viewer", "/api/v1/namespaces/*/pods/*", "GET"},
		{"viewer", "/api/v1/clusters/*", "GET"},
	}
	perms := logicalPermissionsFromPolicies(policies)
	assert.Contains(t, perms, "read:pods")
	assert.Contains(t, perms, "read:clusters")
	assert.NotContains(t, perms, "write:pods")
}

func TestPoliciesForLogicalPermission(t *testing.T) {
	policies, ok := policiesForLogicalPermission("admin:users")
	assert.True(t, ok)
	assert.NotEmpty(t, policies)

	_, ok = policiesForLogicalPermission("not:a:real:permission")
	assert.False(t, ok)
}
