package auth

import (
	"net/http"
	"sync"
	"time"

	"github.com/ciliverse/cilikube/configs"
	"github.com/gin-gonic/gin"
)

// RateLimiter represents a rate limiter for API requests
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	config   *configs.RateLimitConfig
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config *configs.RateLimitConfig) *RateLimiter {
	limiter := &RateLimiter{
		requests: make(map[string][]time.Time),
		config:   config,
	}

	// Start cleanup goroutine
	go limiter.cleanup()

	return limiter
}

// IsAllowed checks if a request from the given IP is allowed
func (rl *RateLimiter) IsAllowed(ip string, requestType string) bool {
	if !rl.config.Enabled {
		return true
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	key := ip + ":" + requestType

	// Get request history for this IP and request type
	requests := rl.requests[key]

	// Determine limits based on request type
	var limit int
	var window time.Duration

	switch requestType {
	case "login":
		limit = rl.config.LoginAttempts
		window = rl.config.LoginWindow
	case "api":
		limit = rl.config.APIRequests
		window = rl.config.APIWindow
	default:
		limit = rl.config.APIRequests
		window = rl.config.APIWindow
	}

	// Remove old requests outside the window
	cutoff := now.Add(-window)
	validRequests := make([]time.Time, 0)
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}

	// Check if we're within the limit
	if len(validRequests) >= limit {
		return false
	}

	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[key] = validRequests

	return true
}

// cleanup removes old entries from the rate limiter
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()

		for key, requests := range rl.requests {
			// Keep only requests from the last hour
			cutoff := now.Add(-time.Hour)
			validRequests := make([]time.Time, 0)

			for _, reqTime := range requests {
				if reqTime.After(cutoff) {
					validRequests = append(validRequests, reqTime)
				}
			}

			if len(validRequests) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = validRequests
			}
		}

		rl.mutex.Unlock()
	}
}

// Global rate limiter instance
var globalRateLimiter *RateLimiter

// InitializeRateLimiter initializes the global rate limiter
func InitializeRateLimiter(config *configs.RateLimitConfig) {
	globalRateLimiter = NewRateLimiter(config)
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(requestType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if globalRateLimiter == nil {
			c.Next()
			return
		}

		ip := c.ClientIP()
		if !globalRateLimiter.IsAllowed(ip, requestType) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// LoginRateLimitMiddleware rate limiting middleware specifically for login attempts
func LoginRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware("login")
}

// APIRateLimitMiddleware rate limiting middleware for general API requests
func APIRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware("api")
}
