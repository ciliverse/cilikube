// pkg/auth/casbin.go
package auth

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath" // 引入 path/filepath

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

// IgnorePath 允许链式调用添加需要忽略的路径
func (r *CasbinBuilder) IgnorePath(path string) *CasbinBuilder {
	r.IgnorePaths = append(r.IgnorePaths, path)
	return r
}

// CasbinMiddleware 返回一个 Gin 中间件处理函数
func (r *CasbinBuilder) CasbinMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqPath := c.Request.URL.Path
		// 跳过 Ignore 的路由
		for _, path := range r.IgnorePaths {
			// 使用 filepath.Match 支持简单的 * 匹配 (如果需要)
			// 或者直接比较 c.Request.URL.Path == path
			if matched, _ := filepath.Match(path, reqPath); matched || reqPath == path {
				c.Next()
				return
			}
		}

		// 从上下文中获取角色 (由 JWT 中间件设置)
		roleVal, exist := c.Get("role")
		if !exist {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无法获取用户角色信息，请先登录"})
			return
		}

		role, ok := roleVal.(string)
		if !ok || role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "用户角色信息格式不正确"})
			return
		}

		obj := reqPath
		act := c.Request.Method

		log.Printf("权限验证 - 角色: %v, 路径: %v, 方法: %v", role, obj, act)

		// 使用 Casbin Enforcer 验证权限
		allowed, err := e.Enforce(role, obj, act)
		if err != nil {
			log.Printf("Casbin Enforce 错误: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "权限检查时发生内部错误"})
			return
		}

		if allowed {
			log.Printf("权限验证通过 - 角色: %s, 路径: %s, 方法: %s", role, obj, act)
			c.Next()
		} else {
			log.Printf("权限验证失败 - 角色: %s 无权访问 %s %s", role, act, obj)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "您没有权限执行此操作"}) // 使用 403 Forbidden
		}
	}
}

// addPolicyIfNotExists 辅助函数，检查策略是否存在，不存在则添加
func addPolicyIfNotExists(e *casbin.Enforcer, sub, obj, act string) {
	has, err := e.HasPolicy(sub, obj, act)
	if err != nil {
		log.Fatalf("检查策略是否存在时出错 (%s, %s, %s): %v", sub, obj, act, err)
	}
	if !has {
		added, err := e.AddPolicy(sub, obj, act)
		if err != nil {
			log.Fatalf("添加策略失败 (%s, %s, %s): %v", sub, obj, act, err)
		}
		if added {
			log.Printf("成功添加默认策略: %s, %s, %s", sub, obj, act)
		} else {
			log.Printf("策略已存在，未添加: %s, %s, %s", sub, obj, act)
		}
	} else {
		log.Printf("策略已存在，跳过添加: %s, %s, %s", sub, obj, act)
	}
}

// InitCasbin 初始化 RBAC 权限控制
func InitCasbin(db *gorm.DB) (*casbin.Enforcer, error) {
	if db == nil {
		return nil, fmt.Errorf("数据库连接 (gorm.DB) 为 nil，无法初始化 Casbin Adapter")
	}

	log.Println("初始化 Casbin Adapter...")
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("创建 Casbin GORM Adapter 失败: %w", err)
	}

	log.Println("初始化 Casbin Enforcer...")
	// 确保 model.conf 路径正确
	e, err := casbin.NewEnforcer("./pkg/auth/model.conf", adapter)
	if err != nil {
		return nil, fmt.Errorf("创建 Casbin Enforcer 失败: %w", err)
	}

	// 启用日志记录 (可选, 但调试时有用)
	e.EnableLog(true)

	// 自动保存策略变更到数据库
	e.EnableAutoSave(true)

	log.Println("从数据库加载策略...")
	if err = e.LoadPolicy(); err != nil {
		log.Printf("加载策略失败 (可能是首次运行，无策略): %v", err)
		// 这里不应该 Fatal，因为首次运行时没有策略是正常的
	}

	log.Println("添加或验证默认策略...")
	// 添加默认权限 (检查是否存在)
	addPolicyIfNotExists(e, "super_admin", "/api/v1/*", "*")   // 管理员拥有所有 v1 接口的所有权限
	addPolicyIfNotExists(e, "normal_user", "/api/v1/*", "GET") // 普通用户只有 GET 权限

	// 你可能还需要添加用户到角色的映射 (g 规则)
	// 例如: e.AddGroupingPolicy("admin", "super_admin")
	// 这通常在用户创建或角色分配时处理，但可以添加默认的。

	// 保存所有可能的新增策略 (如果 AutoSave 不够可靠或需要批量添加)
	// if err := e.SavePolicy(); err != nil {
	//     log.Fatalf("保存策略失败: %v", err)
	// }

	log.Printf("初始化 RBAC 权限控制完成！")
	return e, nil
}
