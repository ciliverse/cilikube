package auth

import "testing"

func TestPathIgnored(t *testing.T) {
	cases := []struct {
		pattern string
		path    string
		want    bool
	}{
		{"/api/v1/auth/login", "/api/v1/auth/login", true},
		{"/api/v1/auth/oauth/*", "/api/v1/auth/oauth/github/auth", true},
		{"/api/v1/auth/oauth/*", "/api/v1/auth/oauth/callback", true},
		{"/api/v1/auth/oauth/*", "/api/v1/auth/login", false},
		{"/api/v1/system/healthz", "/api/v1/system/healthz", true},
	}

	for _, tc := range cases {
		if got := pathIgnored(tc.pattern, tc.path); got != tc.want {
			t.Fatalf("pathIgnored(%q, %q)=%v want %v", tc.pattern, tc.path, got, tc.want)
		}
	}
}
