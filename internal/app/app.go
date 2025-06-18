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
	"github.com/ciliverse/cilikube/internal/initialization"
	"github.com/ciliverse/cilikube/internal/logger"
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

func New(configPath string) (*Application, error) {
	// --- 1. 加载配置 ---
	cfg, err := configs.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	// --- 2. 初始化日志 ---
	var logLevel slog.Level
	if cfg.Server.Mode == "debug" {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}

	appLogger := logger.New(logLevel)
	slog.SetDefault(appLogger)

	// --- 3. 加载配置 ---
	slog.Info("配置加载成功", "path", configPath)

	// --- 4. 数据库和 Store 初始化 ---
	var clusterStore store.ClusterStore
	if cfg.Database.Enabled {
		slog.Info("数据库已启用，正在初始化...")
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
		slog.Info("数据库和 Cluster Store 初始化成功")
	} else {
		slog.Warn("数据库未启用，跳过相关初始化。ClusterStore 将为 nil")
	}

	// --- 5. 初始化 ClusterManager ---
	k8sManager, err := k8s.NewClusterManager(clusterStore, cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化 Kubernetes 集群管理器失败: %w", err)
	}
	slog.Info("Kubernetes 集群管理器初始化成功")

	// --- 6. 初始化应用服务 ---
	services := initialization.InitializeServices(k8sManager, cfg)

	// --- 7. Casbin 初始化 ---
	var e *casbin.Enforcer
	if cfg.Database.Enabled {
		var casbinErr error
		e, casbinErr = auth.InitCasbin(database.DB)
		if casbinErr != nil {
			return nil, fmt.Errorf("初始化 Casbin 失败: %w", casbinErr)
		}
		slog.Info("Casbin 初始化成功")
	}

	// --- 8. Gin 路由器设置 ---
	router := initialization.SetupRouter(cfg, services, k8sManager, e)
	slog.Info("Gin 路由器设置完成")

	return &Application{
		Config: cfg,
		Logger: appLogger,
		Router: router,
	}, nil
}

func (app *Application) Run() {
	serverAddr := ":" + app.Config.Server.Port
	initialization.DisplayServerInfo(serverAddr, app.Config.Server.Mode)
	app.Server = &http.Server{
		Addr:         serverAddr,
		Handler:      app.Router,
		ReadTimeout:  time.Duration(app.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(app.Config.Server.WriteTimeout) * time.Second,
	}
	go func() {
		app.Logger.Info("服务器正在监听...", "address", app.Server.Addr)
		if err := app.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Logger.Error("服务器意外关闭", "error", err)
			os.Exit(1)
		}
	}()
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
