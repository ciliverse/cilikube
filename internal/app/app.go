package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	"github.com/ciliverse/cilikube/configs"
	// 确保引用了 initialization 包，以便调用 displayServerInfo
	"github.com/ciliverse/cilikube/internal/initialization"
	"github.com/ciliverse/cilikube/internal/store"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/ciliverse/cilikube/pkg/database"
	"github.com/ciliverse/cilikube/pkg/k8s"
)

// Application 结构体保持不变
type Application struct {
	Config *configs.Config
	Logger *slog.Logger
	Router *gin.Engine
	Server *http.Server
}

// New 函数保持不变 (为了完整性，这里省略了代码，它和上次一样)
func New(configPath string) (*Application, error) {
	// ... 此处代码与上次重构完全相同 ...
	// --- 1. 初始化日志 ---
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// --- 2. 加载配置 ---
	cfg, err := configs.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}
	logger.Info("配置加载成功", "path", configPath)

	// 根据配置设置日志级别
	var logLevel slog.Level
	if cfg.Server.Mode == "debug" {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger) // 设置为全局默认 logger

	// --- 3. 数据库和 Store 初始化 ---
	var clusterStore store.ClusterStore
	if cfg.Database.Enabled {
		logger.Info("数据库已启用，正在初始化...")
		if err := database.InitDatabase(); err != nil {
			return nil, fmt.Errorf("数据库连接失败: %w", err)
		}
		if err := database.AutoMigrate(); err != nil {
			return nil, fmt.Errorf("数据库自动迁移失败: %w", err)
		}
		if err := database.CreateDefaultAdmin(); err != nil {
			return nil, fmt.Errorf("创建默认管理员失败: %w", err)
		}

		if cfg.Server.EncryptionKey == "" {
			return nil, errors.New("数据库已启用，但配置中未设置 EncryptionKey")
		}
		encryptionKey := []byte(cfg.Server.EncryptionKey)
		clusterStore, err = store.NewGormClusterStore(database.DB, encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("初始化 Cluster Store 失败: %w", err)
		}
		logger.Info("数据库和 Cluster Store 初始化成功")
	} else {
		logger.Warn("数据库未启用，跳过相关初始化。ClusterStore 将为 nil")
	}

	// --- 4. 初始化 ClusterManager ---
	k8sManager, err := k8s.NewClusterManager(clusterStore, cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化 Kubernetes 集群管理器失败: %w", err)
	}
	logger.Info("Kubernetes 集群管理器初始化成功")

	// --- 5. 初始化应用服务和处理器 ---
	services := initialization.InitializeServices(k8sManager, cfg)
	appHandlers := initialization.InitializeHandlers(services, k8sManager)

	// --- 6. Casbin 初始化 ---
	var e *casbin.Enforcer
	if cfg.Database.Enabled {
		var casbinErr error
		e, casbinErr = auth.InitCasbin(database.DB)
		if casbinErr != nil {
			return nil, fmt.Errorf("初始化 Casbin 失败: %w", casbinErr)
		}
		logger.Info("Casbin 初始化成功")
	}

	// --- 7. Gin 路由器设置 ---
	router := initialization.SetupRouter(cfg, appHandlers, e)
	logger.Info("Gin 路由器设置完成")

	return &Application{
		Config: cfg,
		Logger: logger,
		Router: router,
	}, nil
}

// Run 方法被修改以包含启动信息展示
func (app *Application) Run() {
	// --- 1. 构造服务器地址 ---
	serverAddr := ":" + app.Config.Server.Port

	// --- 2. 展示漂亮的启动信息 ---
	// 我们将调用 initialization 包中的公共函数来显示信息。
	// 这要求 displayServerInfo 必须是可导出的（首字母大写）。
	initialization.DisplayServerInfo(serverAddr, app.Config.Server.Mode)

	// --- 3. 配置并启动服务器 ---
	app.Server = &http.Server{
		Addr:         serverAddr,
		Handler:      app.Router,
		ReadTimeout:  time.Duration(app.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(app.Config.Server.WriteTimeout) * time.Second,
	}

	// 启动服务器（goroutine 中）
	go func() {
		// 使用 slog 日志记录启动信息，这对于文件日志和后续处理很有用
		app.Logger.Info("服务器正在监听...", "address", app.Server.Addr)
		if err := app.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Logger.Error("服务器意外关闭", "error", err)
			os.Exit(1)
		}
	}()

	// --- 4. 设置优雅关闭逻辑 ---
	// (这部分代码保持不变)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.Logger.Info("收到关闭信号，正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if app.Config.Database.Enabled {
		database.CloseDatabase()
		app.Logger.Info("数据库连接已关闭")
	}

	if err := app.Server.Shutdown(ctx); err != nil {
		app.Logger.Error("服务器关闭失败", "error", err)
		os.Exit(1)
	}

	app.Logger.Info("服务器已优雅关闭")
}
