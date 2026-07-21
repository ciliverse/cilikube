package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// APIAuditLogger records API audit events without importing internal packages.
type APIAuditLogger interface {
	LogAPIRequest(userID *uint, username, ip, userAgent, resource, action, result, severity string, details map[string]interface{}) error
}

// AuditMiddleware creates middleware for auditing API requests
func AuditMiddleware(auditLogger APIAuditLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if auditLogger == nil || shouldSkipAudit(c.Request.URL.Path) {
			c.Next()
			return
		}

		startTime := time.Now()

		var requestBody []byte
		if c.Request.Body != nil && shouldCaptureBody(c.Request.Method) {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Process request first so JWT middleware can populate user context
		c.Next()

		userID, username, _, hasAuth := GetCurrentUser(c)
		duration := time.Since(startTime)

		details := map[string]interface{}{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"query":       c.Request.URL.RawQuery,
			"status_code": c.Writer.Status(),
			"duration_ms": duration.Milliseconds(),
		}

		if shouldCaptureBody(c.Request.Method) && len(requestBody) > 0 && len(requestBody) < 1024 {
			var bodyMap map[string]interface{}
			if err := json.Unmarshal(requestBody, &bodyMap); err == nil {
				sanitizeRequestBody(bodyMap)
				details["request_body"] = bodyMap
			}
		}

		action, resource := parsePathForAudit(c.Request.Method, c.Request.URL.Path)
		success := c.Writer.Status() < 400

		if !shouldLogEvent(c.Request.Method, c.Request.URL.Path, c.Writer.Status()) {
			return
		}

		_ = auditLogger.LogAPIRequest(
			getUserIDForAudit(userID, hasAuth),
			username,
			c.ClientIP(),
			c.GetHeader("User-Agent"),
			resource,
			action,
			getResultFromStatus(success),
			getSeverityFromStatus(c.Writer.Status()),
			details,
		)
	}
}

func shouldSkipAudit(path string) bool {
	skipPaths := []string{
		"/health",
		"/metrics",
		"/favicon.ico",
		"/static/",
		"/assets/",
		"/uploads/",
	}

	for _, skipPath := range skipPaths {
		if len(path) >= len(skipPath) && path[:len(skipPath)] == skipPath {
			return true
		}
	}
	return false
}

func shouldCaptureBody(method string) bool {
	return method == "POST" || method == "PUT" || method == "PATCH"
}

func sanitizeRequestBody(body map[string]interface{}) {
	sensitiveFields := []string{
		"password", "token", "secret", "key", "auth",
		"credential", "private", "confidential",
	}

	for _, field := range sensitiveFields {
		if _, exists := body[field]; exists {
			body[field] = "[REDACTED]"
		}
	}

	for _, value := range body {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			sanitizeRequestBody(nestedMap)
		}
	}
}

func parsePathForAudit(method, path string) (string, string) {
	resource := "unknown"
	action := method

	if len(path) > 1 {
		parts := splitPath(path)
		if len(parts) >= 3 && parts[0] == "api" {
			if len(parts) >= 3 {
				resource = parts[2] // /api/v1/<resource>
			}
		}
	}

	switch method {
	case "GET":
		action = "read"
	case "POST":
		action = "create"
	case "PUT", "PATCH":
		action = "update"
	case "DELETE":
		action = "delete"
	}

	return action, resource
}

func splitPath(path string) []string {
	var parts []string
	current := ""

	for _, char := range path {
		if char == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}

func getSeverityFromStatus(statusCode int) string {
	switch {
	case statusCode >= 500:
		return "error"
	case statusCode >= 400:
		return "warning"
	default:
		return "info"
	}
}

func getResultFromStatus(success bool) string {
	if success {
		return "success"
	}
	return "failure"
}

func getUserIDForAudit(userID uint, hasAuth bool) *uint {
	if hasAuth && userID > 0 {
		return &userID
	}
	return nil
}

func shouldLogEvent(method, path string, statusCode int) bool {
	if containsString(path, []string{"/auth/", "/login", "/logout"}) {
		return true
	}
	if containsString(path, []string{"/admin/", "/users/", "/roles/", "/settings/"}) {
		return true
	}
	if statusCode >= 400 {
		return true
	}
	if method == "POST" || method == "PUT" || method == "PATCH" || method == "DELETE" {
		return true
	}
	return false
}

func containsString(str string, substrings []string) bool {
	for _, substring := range substrings {
		if len(str) >= len(substring) {
			for i := 0; i <= len(str)-len(substring); i++ {
				if str[i:i+len(substring)] == substring {
					return true
				}
			}
		}
	}
	return false
}
