// pkg/auth/casbin.go
package auth

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath" // Import path/filepath

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
			// Use filepath.Match to support simple * matching (if needed)
			// Or directly compare c.Request.URL.Path == path
			if matched, _ := filepath.Match(path, reqPath); matched || reqPath == path {
				c.Next()
				return
			}
		}

		// Get role from context (set by JWT middleware)
		roleVal, exist := c.Get("role")
		if !exist {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unable to get user role information, please login first"})
			return
		}

		role, ok := roleVal.(string)
		if !ok || role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User role information format is incorrect"})
			return
		}

		obj := reqPath
		act := c.Request.Method

		log.Printf("Permission verification - Role: %v, Path: %v, Method: %v", role, obj, act)

		// Use Casbin Enforcer to verify permissions
		allowed, err := e.Enforce(role, obj, act)
		if err != nil {
			log.Printf("Casbin Enforce error: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal error occurred during permission check"})
			return
		}

		if allowed {
			log.Printf("Permission verification passed - Role: %s, Path: %s, Method: %s", role, obj, act)
			c.Next()
		} else {
			log.Printf("Permission verification failed - Role: %s has no access to %s %s", role, act, obj)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You do not have permission to perform this operation"}) // Use 403 Forbidden
		}
	}
}

// addPolicyIfNotExists helper function, checks if policy exists, adds if not
func addPolicyIfNotExists(e *casbin.Enforcer, sub, obj, act string) {
	has, err := e.HasPolicy(sub, obj, act)
	if err != nil {
		log.Fatalf("Error checking if policy exists (%s, %s, %s): %v", sub, obj, act, err)
	}
	if !has {
		added, err := e.AddPolicy(sub, obj, act)
		if err != nil {
			log.Fatalf("Failed to add policy (%s, %s, %s): %v", sub, obj, act, err)
		}
		if added {
			log.Printf("Successfully added default policy: %s, %s, %s", sub, obj, act)
		} else {
			log.Printf("Policy already exists, not added: %s, %s, %s", sub, obj, act)
		}
	} else {
		log.Printf("Policy already exists, skipping addition: %s, %s, %s", sub, obj, act)
	}
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
	// Ensure model.conf path is correct
	e, err := casbin.NewEnforcer("./pkg/auth/model.conf", adapter)
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

	log.Println("Adding or verifying default policies...")
	// Add default permissions (check if exists)
	addPolicyIfNotExists(e, "super_admin", "/api/v1/*", "*")   // Admin has all permissions for all v1 interfaces
	addPolicyIfNotExists(e, "normal_user", "/api/v1/*", "GET") // Normal users only have GET permissions

	// You may also need to add user to role mappings (g rules)
	// For example: e.AddGroupingPolicy("admin", "super_admin")
	// This is usually handled during user creation or role assignment, but defaults can be added.

	// Save all possible new policies (if AutoSave is not reliable enough or batch addition is needed)
	// if err := e.SavePolicy(); err != nil {
	//     log.Fatalf("Failed to save policies: %v", err)
	// }

	log.Printf("RBAC permission control initialization completed!")
	return e, nil
}
