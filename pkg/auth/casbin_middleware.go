// pkg/auth/casbin.go
package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CasbinBuilder struct {
	IgnorePaths []string
}

func NewCasbinBuilder() *CasbinBuilder {
	return &CasbinBuilder{}
}

// IgnorePath allows chained calls to add paths that need to be ignored
func (r *CasbinBuilder) IgnorePath(path string) *CasbinBuilder {
	r.IgnorePaths = append(r.IgnorePaths, path)
	return r
}

// CasbinMiddleware returns a Gin middleware handler function
func (r *CasbinBuilder) CasbinMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqPath := c.Request.URL.Path
		// Skip ignored routes
		for _, path := range r.IgnorePaths {
			if pathIgnored(path, reqPath) {
				c.Next()
				return
			}
		}

		// Get user ID from context (set by JWT middleware)
		userIDVal, exist := c.Get("user_id")
		if !exist {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unable to get user information, please login first"})
			return
		}

		userID, ok := userIDVal.(uint)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User information format is incorrect"})
			return
		}

		obj := reqPath
		act := c.Request.Method

		log.Printf("Permission verification - UserID: %v, Path: %v, Method: %v", userID, obj, act)

		// Use user-based permission checking instead of role-based
		userSubject := fmt.Sprintf("user:%d", userID)

		// Use Casbin Enforcer to verify permissions
		allowed, err := e.Enforce(userSubject, obj, act)
		if err != nil {
			log.Printf("Casbin Enforce error: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal error occurred during permission check"})
			return
		}

		if allowed {
			log.Printf("Permission verification passed - UserID: %d, Path: %s, Method: %s", userID, obj, act)
			c.Next()
		} else {
			log.Printf("Permission verification failed - UserID: %d has no access to %s %s", userID, act, obj)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You do not have permission to perform this operation"}) // Use 403 Forbidden
		}
	}
}

func pathIgnored(pattern, reqPath string) bool {
	if reqPath == pattern {
		return true
	}
	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		return reqPath == prefix || strings.HasPrefix(reqPath, prefix+"/")
	}
	if strings.HasSuffix(pattern, "/") {
		return strings.HasPrefix(reqPath, pattern)
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(reqPath, strings.TrimSuffix(pattern, "*"))
	}
	matched, _ := filepath.Match(pattern, reqPath)
	return matched
}

// InitCasbin initializes RBAC permission control
func InitCasbin(db *gorm.DB) (*casbin.Enforcer, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection (gorm.DB) is nil, cannot initialize Casbin Adapter")
	}

	log.Println("Initializing Casbin Adapter...")
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create Casbin GORM Adapter: %w", err)
	}

	log.Println("Initializing Casbin Enforcer...")
	modelPath := os.Getenv("CILIKUBE_CASBIN_MODEL")
	if modelPath == "" {
		modelPath = "./pkg/auth/model.conf"
	}
	e, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create Casbin Enforcer: %w", err)
	}

	// Enable logging (optional, but useful for debugging)
	e.EnableLog(true)

	// Auto-save policy changes to database
	e.EnableAutoSave(true)

	log.Println("Loading policies from database...")
	if err = e.LoadPolicy(); err != nil {
		log.Printf("Failed to load policies (may be first run, no policies): %v", err)
		// Should not Fatal here, as having no policies on first run is normal
	}

	// Legacy demo roles — do NOT grant GET /api/v1/* (would allow exec/secrets via keyMatch).
	// Real permissions are seeded by PermissionService.InitializeDefaultPolicies (admin/editor/viewer).
	if _, err := e.RemoveFilteredPolicy(0, "normal_user"); err != nil {
		log.Printf("warning: failed to remove legacy normal_user policies: %v", err)
	}

	log.Printf("RBAC permission control initialization completed!")
	return e, nil
}
