package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditMiddleware creates middleware for auditing API requests
func AuditMiddleware(auditService interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip audit for health checks and static files
		if shouldSkipAudit(c.Request.URL.Path) {
			c.Next()
			return
		}

		startTime := time.Now()

		// Capture request body for audit (if needed)
		var requestBody []byte
		if c.Request.Body != nil && shouldCaptureBody(c.Request.Method) {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Get user information if available
		userID, username, role, hasAuth := GetCurrentUser(c)

		// Process request
		c.Next()

		// Record audit log after request completion
		duration := time.Since(startTime)

		// Prepare audit details
		details := map[string]interface{}{
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"query":        c.Request.URL.RawQuery,
			"status_code":  c.Writer.Status(),
			"duration_ms":  duration.Milliseconds(),
			"content_type": c.GetHeader("Content-Type"),
			"user_agent":   c.GetHeader("User-Agent"),
			"referer":      c.GetHeader("Referer"),
		}

		// Add request body to details if it's a sensitive operation
		if shouldCaptureBody(c.Request.Method) && len(requestBody) > 0 && len(requestBody) < 1024 {
			// Only capture small request bodies and sanitize sensitive data
			var bodyMap map[string]interface{}
			if err := json.Unmarshal(requestBody, &bodyMap); err == nil {
				sanitizeRequestBody(bodyMap)
				details["request_body"] = bodyMap
			}
		}

		// Determine action and resource from path
		action, resource := parsePathForAudit(c.Request.Method, c.Request.URL.Path)

		// Determine if this was a successful operation
		success := c.Writer.Status() < 400

		// Create audit event
		auditEvent := map[string]interface{}{
			"type":       "api_request",
			"severity":   getSeverityFromStatus(c.Writer.Status()),
			"user_id":    getUserIDForAudit(userID, hasAuth),
			"username":   username,
			"user_role":  role,
			"ip_address": c.ClientIP(),
			"user_agent": c.GetHeader("User-Agent"),
			"resource":   resource,
			"action":     action,
			"result":     getResultFromStatus(success),
			"details":    details,
			"timestamp":  startTime,
		}

		// Log the audit event (in a real implementation, you would call the audit service)
		// For now, we'll just log high-priority events
		if shouldLogEvent(c.Request.Method, c.Request.URL.Path, c.Writer.Status()) {
			logAuditEvent(auditEvent)
		}
	}
}

// shouldSkipAudit determines if a request should be skipped from auditing
func shouldSkipAudit(path string) bool {
	skipPaths := []string{
		"/health",
		"/metrics",
		"/favicon.ico",
		"/static/",
		"/assets/",
	}

	for _, skipPath := range skipPaths {
		if len(path) >= len(skipPath) && path[:len(skipPath)] == skipPath {
			return true
		}
	}

	return false
}

// shouldCaptureBody determines if request body should be captured
func shouldCaptureBody(method string) bool {
	return method == "POST" || method == "PUT" || method == "PATCH"
}

// sanitizeRequestBody removes sensitive information from request body
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

	// Also check for nested objects
	for _, value := range body {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			sanitizeRequestBody(nestedMap)
		}
	}
}

// parsePathForAudit extracts action and resource from HTTP method and path
func parsePathForAudit(method, path string) (string, string) {
	// Simple parsing logic - in a real implementation, this would be more sophisticated
	resource := "unknown"
	action := method

	// Extract resource from path
	if len(path) > 1 {
		parts := splitPath(path)
		if len(parts) >= 3 && parts[1] == "api" {
			if len(parts) >= 4 {
				resource = parts[3] // e.g., /api/v1/users -> users
			}
		}
	}

	// Map HTTP methods to actions
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

// splitPath splits a path into components
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

// getSeverityFromStatus determines severity based on HTTP status code
func getSeverityFromStatus(statusCode int) string {
	switch {
	case statusCode >= 500:
		return "error"
	case statusCode >= 400:
		return "warning"
	case statusCode >= 300:
		return "info"
	default:
		return "info"
	}
}

// getResultFromStatus determines result based on success
func getResultFromStatus(success bool) string {
	if success {
		return "success"
	}
	return "failure"
}

// getUserIDForAudit returns user ID for audit or nil if not authenticated
func getUserIDForAudit(userID uint, hasAuth bool) *uint {
	if hasAuth && userID > 0 {
		return &userID
	}
	return nil
}

// shouldLogEvent determines if an event should be logged based on priority
func shouldLogEvent(method, path string, statusCode int) bool {
	// Always log authentication-related requests
	if containsString(path, []string{"/auth/", "/login", "/logout"}) {
		return true
	}

	// Always log admin operations
	if containsString(path, []string{"/admin/", "/users/", "/roles/"}) {
		return true
	}

	// Always log errors
	if statusCode >= 400 {
		return true
	}

	// Log write operations
	if method == "POST" || method == "PUT" || method == "PATCH" || method == "DELETE" {
		return true
	}

	return false
}

// containsString checks if a string contains any of the substrings
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

// logAuditEvent logs an audit event (placeholder implementation)
func logAuditEvent(event map[string]interface{}) {
	// In a real implementation, this would call the audit service
	// For now, we'll just print important events
	if event["severity"] == "error" || event["severity"] == "warning" {
		// This would be replaced with actual audit service call
		// auditService.LogSecurityEvent(event)
	}
}
