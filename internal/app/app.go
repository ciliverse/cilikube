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

type Application struct {
	Config *configs.Config
	Logger *slog.Logger
	Router *gin.Engine
	Server *http.Server
}

func New(configPath string) (*Application, error) {
	// --- 1. Load configuration ---
	cfg, err := configs.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// --- 2. Initialize logger ---
	var logLevel slog.Level
	if cfg.Server.Mode == "debug" {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}

	appLogger := logger.New(logLevel)
	slog.SetDefault(appLogger)

	// --- 3. Configuration loaded ---
	slog.Info("configuration loaded successfully", "path", configPath)

	// --- 4. Database and Store initialization ---
	var clusterStore store.ClusterStore
	if cfg.Database.Enabled {
		slog.Info("database enabled, initializing...")
		if err := database.InitDatabase(); err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		if err := database.AutoMigrate(); err != nil {
			return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
		}
		if err := database.CreateDefaultAdmin(); err != nil {
			return nil, fmt.Errorf("failed to create default admin: %w", err)
		}

		if cfg.Server.EncryptionKey == "" {
			return nil, errors.New("database is enabled but EncryptionKey is not set in configuration")
		}
		encryptionKey := []byte(cfg.Server.EncryptionKey)
		clusterStore, err = store.NewGormClusterStore(database.DB, encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Cluster Store: %w", err)
		}
		slog.Info("database and Cluster Store initialized successfully")
	} else {
		slog.Warn("database not enabled, skipping related initialization. ClusterStore will be nil")
	}

	// --- 5. Initialize ClusterManager ---
	k8sManager, err := k8s.NewClusterManager(clusterStore, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kubernetes cluster manager: %w", err)
	}
	slog.Info("Kubernetes cluster manager initialized successfully")

	// --- 6. Initialize application services ---
	services := initialization.InitializeServices(k8sManager, cfg)

	// --- 7. Casbin initialization ---
	var e *casbin.Enforcer
	if cfg.Database.Enabled {
		var casbinErr error
		e, casbinErr = auth.InitCasbin(database.DB)
		if casbinErr != nil {
			return nil, fmt.Errorf("failed to initialize Casbin: %w", casbinErr)
		}
		slog.Info("Casbin initialized successfully")
	}

	// --- 8. Gin router setup ---
	router := initialization.SetupRouter(cfg, services, k8sManager, e)
	slog.Info("Gin router setup completed")

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
		app.Logger.Info("server is listening...", "address", app.Server.Addr)
		if err := app.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Logger.Error("server closed unexpectedly", "error", err)
			os.Exit(1)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.Logger.Info("received shutdown signal, shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if app.Config.Database.Enabled {
		database.CloseDatabase()
		app.Logger.Info("database connection closed")
	}
	if err := app.Server.Shutdown(ctx); err != nil {
		app.Logger.Error("failed to shutdown server", "error", err)
		os.Exit(1)
	}
	app.Logger.Info("server shutdown gracefully")
}
