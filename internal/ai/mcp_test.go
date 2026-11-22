package ai

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestMCPListIsReadOnlyAndScanCallWorks(t *testing.T) {
	listed, rpcErr := dispatchMCP(context.Background(), rpcRequest{
		JSONRPC: "2.0",
		ID:      json.RawMessage("1"),
		Method:  "tools/list",
	}, nil, nil, "en")
	if rpcErr != nil {
		t.Fatal(rpcErr.Message)
	}
	raw, err := json.Marshal(listed)
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw)
	for _, name := range []string{"scan_health", "get_metrics", "query_prometheus", "get_pod_logs"} {
		if !strings.Contains(text, name) {
			t.Errorf("tools/list missing %s", name)
		}
	}
	for _, banned := range []string{`"name":"delete`, `"name":"create_`, `"name":"scale_`, `"name":"exec`} {
		if strings.Contains(text, banned) {
			t.Fatalf("read-only MCP listed a write tool: %s", banned)
		}
	}

	called, rpcErr := dispatchMCP(context.Background(), rpcRequest{
		JSONRPC: "2.0",
		ID:      json.RawMessage("2"),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"scan_health","arguments":{}}`),
	}, prodShowcaseClient(), nil, "en")
	if rpcErr != nil {
		t.Fatal(rpcErr.Message)
	}
	body, err := json.Marshal(called)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "CrashLoopBackOff") {
		t.Fatalf("mcp scan missing crashloop: %s", body)
	}
	if strings.Contains(string(body), `"isError":true`) {
		t.Fatalf("scan returned isError: %s", body)
	}

	_, rpcErr = dispatchMCP(context.Background(), rpcRequest{Method: "resources/write"}, nil, nil, "en")
	if rpcErr == nil || rpcErr.Code != -32601 {
		t.Fatalf("expected method not found, got %#v", rpcErr)
	}
}
