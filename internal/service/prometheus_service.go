package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ciliverse/cilikube/configs"
)

// PrometheusService queries an external Prometheus server via its HTTP API.
type PrometheusService struct {
	config *configs.Config
	client *http.Client
}

func NewPrometheusService(config *configs.Config) *PrometheusService {
	timeout := 15 * time.Second
	if config != nil && config.Prometheus.Timeout > 0 {
		timeout = config.Prometheus.Timeout
	}
	return &PrometheusService{
		config: config,
		client: &http.Client{Timeout: timeout},
	}
}

func (s *PrometheusService) Enabled() bool {
	return s.config != nil && s.config.Prometheus.Enabled && s.config.Prometheus.URL != ""
}

func (s *PrometheusService) baseURL() string {
	if s.config == nil {
		return ""
	}
	return s.config.Prometheus.URL
}

// PromQueryResult is a simplified Prometheus API response payload.
type PromQueryResult struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
	Error  string          `json:"error,omitempty"`
	Type   string          `json:"errorType,omitempty"`
}

// Query executes an instant PromQL query.
func (s *PrometheusService) Query(ctx context.Context, query string, ts *time.Time) (*PromQueryResult, error) {
	if !s.Enabled() {
		return nil, fmt.Errorf("prometheus integration is disabled or URL is not configured")
	}
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	params := url.Values{}
	params.Set("query", query)
	if ts != nil {
		params.Set("time", strconv.FormatInt(ts.Unix(), 10))
	}
	return s.doGet(ctx, "/api/v1/query", params)
}

// QueryRange executes a range PromQL query.
func (s *PrometheusService) QueryRange(ctx context.Context, query string, start, end time.Time, step string) (*PromQueryResult, error) {
	if !s.Enabled() {
		return nil, fmt.Errorf("prometheus integration is disabled or URL is not configured")
	}
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}
	if step == "" {
		step = "60s"
	}

	params := url.Values{}
	params.Set("query", query)
	params.Set("start", strconv.FormatInt(start.Unix(), 10))
	params.Set("end", strconv.FormatInt(end.Unix(), 10))
	params.Set("step", step)
	return s.doGet(ctx, "/api/v1/query_range", params)
}

// GetStatus returns basic connectivity info for the configured Prometheus.
func (s *PrometheusService) GetStatus(ctx context.Context) (map[string]interface{}, error) {
	if !s.Enabled() {
		return map[string]interface{}{
			"enabled": false,
			"url":     "",
			"healthy": false,
		}, nil
	}

	result, err := s.doGet(ctx, "/api/v1/status/buildinfo", nil)
	if err != nil {
		return map[string]interface{}{
			"enabled": true,
			"url":     s.baseURL(),
			"healthy": false,
			"error":   err.Error(),
		}, nil
	}

	return map[string]interface{}{
		"enabled": true,
		"url":     s.baseURL(),
		"healthy": result.Status == "success",
		"data":    result.Data,
	}, nil
}

func (s *PrometheusService) doGet(ctx context.Context, path string, params url.Values) (*PromQueryResult, error) {
	base := strings.TrimRight(s.baseURL(), "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	endpoint := base + path
	if len(params) > 0 {
		endpoint = endpoint + "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("prometheus request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, fmt.Errorf("failed to read prometheus response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("prometheus returned status %d: %s", resp.StatusCode, string(body))
	}

	var result PromQueryResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode prometheus response: %w", err)
	}
	if result.Status != "success" && result.Error != "" {
		return &result, fmt.Errorf("prometheus query error: %s", result.Error)
	}
	return &result, nil
}
