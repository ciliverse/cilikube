package auth

import "testing"

func TestShouldLogEventIncludesAPIReads(t *testing.T) {
	cases := []struct {
		method string
		path   string
		status int
		want   bool
	}{
		{"GET", "/api/v1/namespaces/default/pods", 200, true},
		{"GET", "/api/v1/clusters", 200, true},
		{"GET", "/api/v1/ai/skills", 200, true},
		{"GET", "/api/v1/fleet/summary", 200, true},
		{"GET", "/api/v1/monitoring/nodes", 200, false}, // noisy
		{"GET", "/api/v1/audit/logs", 200, false},       // self
		{"GET", "/api/v1/namespaces/default/pods/x/watch", 200, false},
		{"GET", "/health", 200, false}, // skipped earlier in middleware; still false here
		{"POST", "/api/v1/namespaces/default/pods", 201, true},
		{"GET", "/api/v1/auth/oauth/providers", 200, true},
		{"GET", "/api/v1/pods", 500, true},
	}
	for _, tc := range cases {
		got := shouldLogEvent(tc.method, tc.path, tc.status)
		if got != tc.want {
			t.Fatalf("%s %s %d: got %v want %v", tc.method, tc.path, tc.status, got, tc.want)
		}
	}
}

func TestShouldDedupeAuditRead(t *testing.T) {
	keyUser := "u-test-dedupe"
	path := "/api/v1/clusters"
	if shouldDedupeAuditRead(keyUser, "GET", path, 200) {
		t.Fatal("first hit should not dedupe")
	}
	if !shouldDedupeAuditRead(keyUser, "GET", path, 200) {
		t.Fatal("second hit within window should dedupe")
	}
	if shouldDedupeAuditRead(keyUser, "GET", "/api/v1/nodes", 200) {
		t.Fatal("different path should not dedupe")
	}
}
