package service

import (
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

func newTestViewerEnforcer(t *testing.T) *casbin.Enforcer {
	t.Helper()
	text := `
[request_definition]
r = sub, obj, act
[policy_definition]
p = sub, obj, act
[role_definition]
g = _, _
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
`
	m, err := model.NewModelFromString(text)
	if err != nil {
		t.Fatalf("model: %v", err)
	}
	e, err := casbin.NewEnforcer(m)
	if err != nil {
		t.Fatalf("enforcer: %v", err)
	}
	for _, p := range systemViewerPolicies() {
		if _, err := e.AddPolicy(p.role, p.object, p.action); err != nil {
			t.Fatalf("add policy: %v", err)
		}
	}
	if _, err := e.AddGroupingPolicy("user:1", "viewer"); err != nil {
		t.Fatalf("grouping: %v", err)
	}
	return e
}

func TestViewerAllowsReadDeniesDangerous(t *testing.T) {
	e := newTestViewerEnforcer(t)
	sub := "user:1"

	allow := []struct{ obj, act string }{
		{"/api/v1/pods", "GET"},
		{"/api/v1/namespaces/default/pods", "GET"},
		{"/api/v1/namespaces/default/pods/nginx", "GET"},
		{"/api/v1/namespaces/default/pods/nginx/logs", "GET"},
		{"/api/v1/pods/metrics", "GET"},
		{"/api/v1/namespaces", "GET"},
		{"/api/v1/namespaces/default", "GET"},
		{"/api/v1/nodes", "GET"},
		{"/api/v1/clusters", "GET"},
		{"/api/v1/configmaps", "GET"},
	}
	for _, tc := range allow {
		ok, err := e.Enforce(sub, tc.obj, tc.act)
		if err != nil {
			t.Fatalf("enforce %s %s: %v", tc.act, tc.obj, err)
		}
		if !ok {
			t.Errorf("viewer should allow %s %s", tc.act, tc.obj)
		}
	}

	deny := []struct{ obj, act string }{
		{"/api/v1/namespaces/default/pods/nginx/exec", "GET"},
		{"/api/v1/namespaces/default/pods/nginx/attach", "GET"},
		{"/api/v1/namespaces/default/pods/nginx/portforward", "GET"},
		{"/api/v1/secrets", "GET"},
		{"/api/v1/namespaces/default/secrets", "GET"},
		{"/api/v1/namespaces/default/secrets/mysecret", "GET"},
		{"/api/v1/proxy/pods", "GET"},
		{"/api/v1/namespaces/default/pods", "POST"},
		{"/api/v1/namespaces/default/pods/nginx", "DELETE"},
		{"/api/v1/namespaces/default/deployments/x", "PUT"},
		{"/api/v1/admin/users", "GET"},
		{"/api/v1/clusters", "POST"},
	}
	for _, tc := range deny {
		ok, err := e.Enforce(sub, tc.obj, tc.act)
		if err != nil {
			t.Fatalf("enforce %s %s: %v", tc.act, tc.obj, err)
		}
		if ok {
			t.Errorf("viewer must deny %s %s", tc.act, tc.obj)
		}
	}
}

func TestNamespacesStarDoesNotMatchNestedWithParamStyle(t *testing.T) {
	// Guard: :ns must not match nested exec paths (the old /namespaces/* bug).
	e := newTestViewerEnforcer(t)
	ok, err := e.Enforce("user:1", "/api/v1/namespaces/default/pods/x/exec", "GET")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("namespaces/:ns must not allow nested exec")
	}
}
