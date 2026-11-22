package ai

import (
	"context"
	"strings"
	"testing"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
)

func prodShowcaseClient() *k8s.Client {
	return k8s.NewShowcaseClientAt(
		k8s.ShowcaseProdAPIHost,
		k8s.ShowcaseProdVersion,
		k8s.ShowcaseProdSeedObjects()...,
	)
}

func TestScanHealthFindsCrashLoopAndNotReadyNode(t *testing.T) {
	res, err := executeTool(context.Background(), prodShowcaseClient(), nil, "scan_health", nil, "en")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"orders-api-a1b2c-55555",
		"CrashLoopBackOff",
		"prod-worker-3",
		"NotReady",
		"inventory-f1e2d-eeee2",
		"Failed",
	} {
		if !strings.Contains(res.Text, want) {
			t.Errorf("scan missing %q\n%s", want, res.Text)
		}
	}
	if len(res.Resources) == 0 {
		t.Fatal("expected resource refs")
	}
}

func TestScanHealthNamespaceSkipsOtherNamespaces(t *testing.T) {
	res, err := executeTool(context.Background(), prodShowcaseClient(), nil, "scan_health", map[string]interface{}{
		"namespace": "payments",
	}, "zh")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Text, "pay-gateway-9x8y7-p3333") {
		t.Fatalf("expected payments crashloop, got %s", res.Text)
	}
	if strings.Contains(res.Text, "orders-api-a1b2c-55555") {
		t.Fatalf("namespace filter leaked production pod:\n%s", res.Text)
	}
	if !strings.Contains(res.Text, "prod-worker-3") {
		t.Fatal("nodes should stay cluster-wide")
	}
}

func TestGetMetricsShowcaseNodes(t *testing.T) {
	res, err := executeTool(context.Background(), prodShowcaseClient(), nil, "get_metrics", map[string]interface{}{
		"scope": "node",
	}, "en")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Text, "prod-worker-1") || !strings.Contains(res.Text, "cpu=") {
		t.Fatalf("expected node metrics, got %s", res.Text)
	}
}

func TestQueryPrometheusSkipsWhenUnset(t *testing.T) {
	res, err := executeTool(context.Background(), prodShowcaseClient(), service.NewPrometheusService(nil), "query_prometheus", map[string]interface{}{
		"query": "up",
	}, "en")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Text, "not configured") {
		t.Fatalf("expected skip, got %s", res.Text)
	}

	_, err = executeTool(context.Background(), prodShowcaseClient(), nil, "query_prometheus", map[string]interface{}{
		"query": strings.Repeat("a", 801),
	}, "en")
	if err == nil || !strings.Contains(err.Error(), "800") {
		t.Fatalf("expected length error, got %v", err)
	}
}
