package auth

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	SessionID string `json:"session_id,omitempty"`
	jwt.RegisteredClaims
}

// SessionValidator validates persisted sessions referenced by JWT claims.
type SessionValidator interface {
	ValidateSessionID(sessionID string) error
}

var sessionValidator SessionValidator

// SetSessionValidator configures optional session validation for JWTAuthMiddleware.
func SetSessionValidator(v SessionValidator) {
	sessionValidator = v
}

// GenerateToken generates JWT token
func GenerateToken(user *models.User, sessionID string) (string, time.Time, error) {
	expirationTime := time.Now().Add(configs.GlobalConfig.JWT.ExpireDuration)

	claims := &JWTClaims{
		UserID:    user.ID,
		Username:  user.Username,
		Role:      user.Role,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    configs.GlobalConfig.JWT.Issuer,
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(configs.GlobalConfig.JWT.SecretKey))

	return tokenString, expirationTime, err
}

// ParseToken parses JWT token
func ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.GlobalConfig.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

// JWTAuthUnless runs JWT auth for all paths except the provided prefixes/patterns.
func JWTAuthUnless(skipPaths ...string) gin.HandlerFunc {
	jwtAuth := JWTAuthMiddleware()
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, skip := range skipPaths {
			if path == skip || strings.HasPrefix(path, skip) {
				c.Next()
				return
			}
			if matched, _ := filepath.Match(skip, path); matched {
				c.Next()
				return
			}
		}
		jwtAuth(c)
	}
}

// extractBearerToken resolves JWT from Authorization header or query (WebSocket clients
// cannot set custom headers, so ?token= is accepted for upgrade requests).
func extractBearerToken(c *gin.Context) (string, string) {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimSpace(authHeader[7:]), ""
	}
	if authHeader != "" {
		return "", "Invalid authorization header format"
	}
	if token := strings.TrimSpace(c.Query("token")); token != "" {
		return token, ""
	}
	if token := strings.TrimSpace(c.Query("access_token")); token != "" {
		return token, ""
	}
	return "", "Authorization header or token query parameter is required"
}

// JWTAuthMiddleware JWT authentication middleware
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, extractErr := extractBearerToken(c)
		if extractErr != "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": extractErr,
			})
			c.Abort()
			return
		}

		// Parse token
		claims, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid token: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Check if token is expired
		if claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Token has expired",
			})
			c.Abort()
			return
		}

		// Validate persisted session when configured (enables logout revocation)
		if sessionValidator != nil {
			if claims.SessionID == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "Session information missing from token",
				})
				c.Abort()
				return
			}
			if err := sessionValidator.ValidateSessionID(claims.SessionID); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "Session invalid: " + err.Error(),
				})
				c.Abort()
				return
			}
		}

		// Store user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)
		c.Set("session_id", claims.SessionID)

		c.Next()
	}
}

// JWTAuthWithSessionMiddleware enables session validation and returns JWT middleware.
func JWTAuthWithSessionMiddleware(validator SessionValidator) gin.HandlerFunc {
	SetSessionValidator(validator)
	return JWTAuthMiddleware()
}

// AdminRequiredMiddleware admin privilege middleware
func AdminRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			fmt.Printf("DEBUG: User role not found in context\n")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "User information not found",
			})
			c.Abort()
			return
		}

		fmt.Printf("DEBUG: User role from JWT: %v (type: %T)\n", role, role)

		if role != "admin" {
			fmt.Printf("DEBUG: Access denied - role '%v' is not admin\n", role)
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "Admin privileges required",
			})
			c.Abort()
			return
		}

		fmt.Printf("DEBUG: Admin access granted for role: %v\n", role)
		c.Next()
	}
}

// OptionalAuthMiddleware optional authentication middleware (does not require mandatory login)
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		} else {
			c.Next()
			return
		}

		claims, err := ParseToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		if claims.ExpiresAt.Time.Before(time.Now()) {
			c.Next()
			return
		}

		// Set user information to context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// GetCurrentUser gets current user information from context
func GetCurrentUser(c *gin.Context) (uint, string, string, bool) {
	userID, exists1 := c.Get("user_id")
	username, exists2 := c.Get("username")
	role, exists3 := c.Get("user_role")

	if !exists1 || !exists2 || !exists3 {
		return 0, "", "", false
	}

	return userID.(uint), username.(string), role.(string), true
}
