package service

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ciliverse/cilikube/pkg/k8s"
)

// showcasePrometheusActive uses in-process synthetic PromQL when the public
// exhibit has no real Prometheus URL configured.
func (s *PrometheusService) showcasePrometheusActive() bool {
	if !k8s.IsShowcase() {
		return false
	}
	if s.config == nil || !s.config.Prometheus.Enabled {
		return true
	}
	u := strings.TrimSpace(s.config.Prometheus.URL)
	if u == "" {
		return true
	}
	lower := strings.ToLower(u)
	if strings.HasPrefix(lower, "showcase:") || strings.Contains(lower, "showcase.cilikube.local") {
		return true
	}
	return false
}

func showcasePromStatus() map[string]interface{} {
	return map[string]interface{}{
		"enabled": true,
		"healthy": true,
		"url":     "showcase://prometheus",
		"mode":    "showcase",
	}
}

func parsePromStep(step string) time.Duration {
	if step == "" {
		return time.Minute
	}
	d, err := time.ParseDuration(step)
	if err != nil || d <= 0 {
		return time.Minute
	}
	if d < 15*time.Second {
		return 15 * time.Second
	}
	return d
}

func showcaseSeriesKind(query string) string {
	q := strings.ToLower(query)
	switch {
	case strings.Contains(q, "cpu"):
		return "cpu"
	case strings.Contains(q, "memory") || strings.Contains(q, "mem"):
		return "memory"
	default:
		return "generic"
	}
}

func showcaseSample(kind string, t time.Time) float64 {
	sec := float64(t.Unix())
	switch kind {
	case "cpu":
		// cores
		v := 1.8 + 0.9*math.Sin(sec/90) + 0.35*math.Sin(sec/37)
		if v < 0.4 {
			v = 0.4
		}
		return v
	case "memory":
		// GiB
		v := 12.5 + 2.8*math.Cos(sec/120) + 1.1*math.Sin(sec/55)
		if v < 4 {
			v = 4
		}
		return v
	default:
		v := 50 + 20*math.Sin(sec/60)
		if v < 5 {
			v = 5
		}
		return v
	}
}

func showcaseQueryRangeResult(query string, start, end time.Time, step string) (*PromQueryResult, error) {
	if !end.After(start) {
		return nil, fmt.Errorf("end must be after start")
	}
	kind := showcaseSeriesKind(query)
	interval := parsePromStep(step)
	// Cap points to keep payloads small.
	maxPoints := 240
	est := int(end.Sub(start)/interval) + 1
	if est > maxPoints {
		interval = end.Sub(start) / time.Duration(maxPoints)
		if interval < 15*time.Second {
			interval = 15 * time.Second
		}
	}

	values := make([][]interface{}, 0, maxPoints+1)
	for ts := start; !ts.After(end); ts = ts.Add(interval) {
		v := showcaseSample(kind, ts)
		values = append(values, []interface{}{ts.Unix(), strconv.FormatFloat(v, 'f', 4, 64)})
		if len(values) >= maxPoints {
			break
		}
	}
	if len(values) == 0 {
		values = append(values, []interface{}{end.Unix(), strconv.FormatFloat(showcaseSample(kind, end), 'f', 4, 64)})
	}

	payload := map[string]interface{}{
		"resultType": "matrix",
		"result": []map[string]interface{}{
			{
				"metric": map[string]string{
					"job":       "showcase",
					"__name__":  "showcase_" + kind,
					"cluster":   "demo",
					"simulated": "true",
				},
				"values": values,
			},
		},
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &PromQueryResult{Status: "success", Data: raw}, nil
}

func showcaseQueryResult(query string, ts *time.Time) (*PromQueryResult, error) {
	at := time.Now()
	if ts != nil {
		at = *ts
	}
	kind := showcaseSeriesKind(query)
	v := showcaseSample(kind, at)
	payload := map[string]interface{}{
		"resultType": "vector",
		"result": []map[string]interface{}{
			{
				"metric": map[string]string{
					"job":      "showcase",
					"__name__": "showcase_" + kind,
				},
				"value": []interface{}{at.Unix(), strconv.FormatFloat(v, 'f', 4, 64)},
			},
		},
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &PromQueryResult{Status: "success", Data: raw}, nil
}
