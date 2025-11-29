package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/configs"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken generates JWT token
func GenerateToken(user *models.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(configs.GlobalConfig.JWT.ExpireDuration)

	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
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

// JWTAuthMiddleware JWT authentication middleware
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check Bearer prefix
		tokenString := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:] // Remove "Bearer " prefix
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid authorization header format",
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

		// Store user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// AdminRequiredMiddleware admin privilege middleware
func AdminRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "User information not found",
			})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "Admin privileges required",
			})
			c.Abort()
			return
		}

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
