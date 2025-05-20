package auth

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type CasbinBuilder struct {
	IgnorePaths []string
}

func NewCasbinBuilder() *CasbinBuilder {
	return &CasbinBuilder{}
}

func (r *CasbinBuilder) IgnorePath(path string) *CasbinBuilder {
	r.IgnorePaths = append(r.IgnorePaths, path)
	return r
}

func (r *CasbinBuilder) CasbinMiddlerware(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过 Ignore 的路由
		for _, path := range r.IgnorePaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}
		// 获取当前用户
		user, exist := c.Get("username")
		if !exist {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "未登录或无有效认证信息",
			})
			return
		}
		obj := c.Request.URL.Path
		act := c.Request.Method

		// 获取用户的角色
		roles, rolesExist := c.Get("roles")
		if !rolesExist {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "未登录或无有效认证信息",
			})
			return
		}

		// 添加调试日志
		log.Printf("权限验证 - 用户: %v, 角色: %v, 路径: %v, 方法: %v", user, roles, obj, act)

		// 验证用户权限
		if ok, _ := e.Enforce(user, obj, act); ok {
			c.Next()
			log.Printf("权限验证通过 - %s用户权限匹配", user)
			return
		}

		// 验证用户角色权限
		if ok, err := e.Enforce(roles, obj, act); ok {
			c.Next()
			log.Printf("权限验证通过 - %s角色权限匹配", roles)
			return
		} else if err != nil {
			log.Printf("权限验证错误: %v", err)
		}
		// 如果所有权限验证都失败
		log.Printf("权限验证失败 - 无访问权限")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "没有权限访问该资源",
		})
		return
	}
}

// InitCasbin 初始化 RBAC 权限控制
func InitCasbin(db *gorm.DB) (*casbin.Enforcer, error) {
	// --- 初始化 adapter ---
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	// --- 初始化 enforcer ---
	e, err := casbin.NewEnforcer("./pkg/auth/model.conf", adapter)
	if err != nil {
		return nil, err
	}
	// 自动保存
	e.EnableAutoSave(true)

	// --- 从数据库中加载策略 ---
	if err = e.LoadPolicy(); err != nil {
		log.Fatalf("加载策略失败: %v", err)
		return e, err
	}

	// --- 添加默认权限 ---
	// ·超级管理员
	if _, err := e.AddPolicy("super_admin", "/api/v1/*", "*"); err != nil {
		log.Fatalf("添加默认权限（管理员）失败: %v", err)
		return e, err
	}
	// ·普通用户
	if _, err := e.AddPolicy("normal_user", "/api/v1/*", "GET"); err != nil {
		log.Fatalf("添加默认权限（普通用户）失败: %v", err)
		return e, err
	}

	log.Printf("初始化 RBAC 权限控制完成！")
	return e, nil
}
