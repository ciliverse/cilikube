package service

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/ciliverse/cilikube/configs"
)

func TestShowcasePrometheusQueryRange(t *testing.T) {
	t.Setenv("CILIKUBE_SHOWCASE", "1")
	defer os.Unsetenv("CILIKUBE_SHOWCASE")

	svc := NewPrometheusService(&configs.Config{
		Prometheus: configs.PrometheusConfig{Enabled: false},
	})
	if !svc.Enabled() {
		t.Fatal("showcase should enable prometheus")
	}
	st, err := svc.GetStatus(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if st["mode"] != "showcase" || st["healthy"] != true {
		t.Fatalf("status %#v", st)
	}

	end := time.Now()
	start := end.Add(-time.Hour)
	res, err := svc.QueryRange(context.Background(),
		`sum(rate(container_cpu_usage_seconds_total{container!=""}[5m]))`,
		start, end, "60s")
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "success" {
		t.Fatalf("status %s", res.Status)
	}
	var data struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Values [][]interface{} `json:"values"`
		} `json:"result"`
	}
	if err := json.Unmarshal(res.Data, &data); err != nil {
		t.Fatal(err)
	}
	if data.ResultType != "matrix" || len(data.Result) == 0 || len(data.Result[0].Values) < 10 {
		t.Fatalf("unexpected matrix %#v", data)
	}
}

func TestRemotePreferredOverShowcase(t *testing.T) {
	t.Setenv("CILIKUBE_SHOWCASE", "1")
	defer os.Unsetenv("CILIKUBE_SHOWCASE")

	svc := NewPrometheusService(&configs.Config{
		Prometheus: configs.PrometheusConfig{
			Enabled: true,
			URL:     "http://127.0.0.1:9",
		},
	})
	if svc.showcasePrometheusActive() {
		t.Fatal("explicit remote URL should win over showcase mock")
	}
}
