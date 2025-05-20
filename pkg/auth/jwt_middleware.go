package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var expireDuration = 24 * time.Hour

type Claims struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
	Roles    string `json:"roles"`
	jwt.RegisteredClaims
}

type JWTBuilder struct {
	IgnorePaths []string
}

func NewJWTBuilder() *JWTBuilder {
	return &JWTBuilder{}
}

func (j *JWTBuilder) IgnorePath(path string) *JWTBuilder {
	j.IgnorePaths = append(j.IgnorePaths, path)
	return j
}

// GenerateJWT 登录的时候生成 JWT Token
func GenerateJWT(username string, userID uint, roles string, secret string) (string, error) {
	// 创建 JWT Token 的过期时间
	expireTime := time.Now().Add(expireDuration)
	claims := &Claims{
		Username: username,
		UserID:   userID,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   username,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// JWTMiddleware 特定路径访问的时候验证 JWT Token
func (j *JWTBuilder) JWTMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过 Ignore 的路由
		for _, path := range j.IgnorePaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		// 从Header提取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "请你登录之后在访问该资源",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			// 将用户信息存入上下文
			c.Set("username", claims.Username)
			c.Set("user_id", claims.UserID)
			c.Set("roles", claims.Roles)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token: " + err.Error(),
			})
		}
	}
}
