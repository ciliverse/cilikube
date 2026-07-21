package routes

import (
	"github.com/ciliverse/cilikube/internal/handlers"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/gin-gonic/gin"
)

// RegisterMonitoringRoutes registers monitoring, prometheus and informer routes.
func RegisterMonitoringRoutes(router *gin.RouterGroup, monitoringService *service.MonitoringService, prometheusService *service.PrometheusService, k8sManager *k8s.ClusterManager) {
	if monitoringService != nil {
		h := handlers.NewMonitoringHandler(monitoringService)
		monitoring := router.Group("/monitoring")
		{
			monitoring.GET("/metrics", h.GetRealTimeMetrics)
			monitoring.GET("/health", h.GetSystemHealth)
			monitoring.GET("/dashboard", h.GetDashboardData)
			monitoring.GET("/security", h.GetSecurityOverview)
			monitoring.GET("/alerts", h.GetAlerts)
			monitoring.GET("/metrics/history", h.GetMetricsHistory)
		}
	}

	if prometheusService != nil {
		ph := handlers.NewPrometheusHandler(prometheusService)
		prom := router.Group("/prometheus")
		{
			prom.GET("/status", ph.GetStatus)
			prom.GET("/query", ph.Query)
			prom.GET("/query_range", ph.QueryRange)
		}
	}

	if k8sManager != nil {
		ih := handlers.NewInformerHandler(k8sManager)
		informers := router.Group("/informers")
		{
			informers.GET("/status", ih.GetStatus)
			informers.GET("/pods", ih.ListPods)
			informers.GET("/nodes", ih.ListNodes)
			informers.GET("/namespaces", ih.ListNamespaces)
			informers.GET("/services", ih.ListServices)
			informers.GET("/deployments", ih.ListDeployments)
		}
	}
}
