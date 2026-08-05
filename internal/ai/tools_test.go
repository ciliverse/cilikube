package ai

import (
	"context"
	"strings"
	"testing"

	"github.com/ciliverse/cilikube/pkg/k8s"
)

func TestListPodsPhasesClusterWideIgnoresDefaultTrap(t *testing.T) {
	client := k8s.NewShowcaseClientAt(
		k8s.ShowcaseStagingAPIHost,
		k8s.ShowcaseStagingVersion,
		k8s.ShowcaseStagingSeedObjects()...,
	)
	// Same trap as the bad screenshot: UI namespace=default while Pending lives in preview.
	res, err := executeTool(context.Background(), client, "list_resources", map[string]interface{}{
		"kind":      "pods",
		"namespace": "default",
		"phases":    "Pending,Failed,Unknown",
		"limit":     20,
	}, "en")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Text, "no items") {
		t.Fatalf("expected empty in default ns, got %q", res.Text)
	}

	res, err = executeTool(context.Background(), client, "list_resources", map[string]interface{}{
		"kind":   "pods",
		"phases": "Pending,Failed,Unknown",
		"limit":  20,
	}, "en")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Text, "feature-x-pr42-eee05") || !strings.Contains(res.Text, "Pending") {
		t.Fatalf("expected pending pod cluster-wide, got %q", res.Text)
	}
}

func TestListWarningEventsClusterWide(t *testing.T) {
	client := k8s.NewShowcaseClientAt(
		k8s.ShowcaseStagingAPIHost,
		k8s.ShowcaseStagingVersion,
		k8s.ShowcaseStagingSeedObjects()...,
	)
	res, err := executeTool(context.Background(), client, "list_resources", map[string]interface{}{
		"kind":        "events",
		"event_types": "Warning,Error",
		"limit":       30,
	}, "en")
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasPrefix(res.Text, "no items") {
		t.Fatalf("expected warning events, got %q", res.Text)
	}
	if !strings.Contains(res.Text, "[Warning]") {
		t.Fatalf("expected Warning lines, got %q", res.Text)
	}
}
