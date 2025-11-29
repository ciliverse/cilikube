package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors handles cross-origin requests, supports preflight requests
// Suggestion: pass allowedOrigins as parameter or load from configuration
func Cors(allowedOrigins []string) gin.HandlerFunc {
	// If allowedOrigins is empty, provide a default value or print warning
	if len(allowedOrigins) == 0 {
		log.Println("CORS Warning: No allowed origins configured!")
		// Can choose to completely disable CORS or allow all (if AllowCredentials is false)
		// Here we choose to disable CORS and just continue processing
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		// Get the request's Origin
		origin := c.Request.Header.Get("Origin")

		// If there's no Origin header (e.g., non-browser requests or same-origin requests), no need to handle CORS
		if origin == "" {
			c.Next()
			return
		}

		// Check if the request's Origin is in the allowed list
		allowedOrigin := ""
		for _, o := range allowedOrigins {
			if o == origin {
				allowedOrigin = origin
				break
			}
			// Optional: handle more complex matching logic like wildcard subdomains
		}

		// If Origin matches successfully
		if allowedOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			// Important: Since Allow-Origin is not "*", we need the Vary header
			c.Header("Vary", "Origin")

			// Set other CORS headers
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")            // More comprehensive method list
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-CSRF-Token, Accept") // Add common headers
			c.Header("Access-Control-Allow-Credentials", "true")                                          // Support credentials
			c.Header("Access-Control-Max-Age", "86400")                                                   // Preflight request cache time

			// Properly handle OPTIONS preflight requests: if Origin is allowed and method is OPTIONS, abort and return 204
			if c.Request.Method == "OPTIONS" {
				log.Printf("CORS: Preflight request for %s from %s allowed.", c.Request.URL.Path, origin)
				c.AbortWithStatus(http.StatusNoContent)
				return
			}

			// For non-OPTIONS requests, continue processing
			log.Printf("CORS: Allowed non-preflight request for %s from %s.", c.Request.URL.Path, origin)
			c.Next()

		} else {
			// If Origin doesn't match
			log.Printf("CORS: Origin '%s' not allowed for %s.", origin, c.Request.URL.Path)

			// For OPTIONS preflight requests, if Origin is not allowed, should also abort but may not set CORS headers
			// Browser will reject the request due to missing necessary Allow-Origin header
			if c.Request.Method == "OPTIONS" {
				// Can choose to return 403 Forbidden or simply Abort
				// c.AbortWithStatus(http.StatusForbidden) // More explicit rejection
				c.Abort() // Or just abort the processing chain
				return
			}

			// For non-OPTIONS requests, Origin not allowed
			// Browser will send the request but will block frontend JS from reading the response.
			// Here we can choose:
			// 1. Call c.Next(): let request continue, but browser will error (current code logic)
			// 2. Call c.AbortWithStatus(http.StatusForbidden): directly reject request (stricter)
			// Choosing c.Next() means backend might execute operations, but frontend can't receive results.
			// Choosing Abort is safer, preventing unauthorized origins from triggering operations.
			// Here we choose the safer Abort:
			log.Printf("CORS: Aborting non-preflight request from disallowed origin '%s' for %s.", origin, c.Request.URL.Path)
			c.AbortWithStatus(http.StatusForbidden) // Direct rejection
			// Or, if you want to keep the original behavior (allow backend processing but browser blocks):
			// c.Next()

			return // Ensure return after abort
		}
	}
}
