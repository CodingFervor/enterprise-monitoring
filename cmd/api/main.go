package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CodingFervor/enterprise-monitoring/internal/cache"
	"github.com/CodingFervor/enterprise-monitoring/internal/config"
	"github.com/CodingFervor/enterprise-monitoring/internal/database"
	"github.com/CodingFervor/enterprise-monitoring/internal/middleware"
	"github.com/CodingFervor/enterprise-monitoring/internal/service"
	"github.com/CodingFervor/enterprise-monitoring/pkg/jwt"
	"github.com/CodingFervor/enterprise-monitoring/pkg/logger"
)

var svc *service.Context

func main() {
	// Load config
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	// Logger
	logger.SetLevel(cfg.Server.Mode)
	logger.Info("starting " + "enterprise-monitoring")

	// Database
	if err := database.Connect(cfg.Database); err != nil {
		logger.Error("database connect failed", "error", err)
		log.Fatalf("database: %v", err)
	}
	defer database.Close()

	// Redis
	if err := cache.Connect(cfg.Redis); err != nil {
		logger.Warn("redis connect failed, running without cache", "error", err)
	}
	defer cache.Close()

	// JWT & Service
	jwt.SetSecret(cfg.JWT.Secret)
	svc = service.NewContext()

	// Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(requestLogger())

	r.GET("/health", healthCheck)
	r.GET("/ready", readinessCheck)

	api := r.Group("/api/v1")
	{

		api.POST("/auth/login", Login)

		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.POST("/traces/ingest", IngestTrace)
			auth.POST("/traces/batch", IngestTraces)
			auth.GET("/traces", ListTraces)
			auth.GET("/traces/:id", GetTrace)
			auth.GET("/traces/:id/spans", ListTraceSpans)
			auth.GET("/traces/services", ListTraceServices)

			auth.POST("/logs/ingest", IngestLog)
			auth.POST("/logs/batch", IngestLogs)
			auth.GET("/logs", QueryLogs)
			auth.GET("/logs/fields", ListLogFields)
			auth.GET("/logs/aggregations", LogAggregations)

			auth.POST("/metrics/ingest", IngestMetric)
			auth.POST("/metrics/query", QueryMetrics)
			auth.GET("/metrics/labels", ListMetricLabels)
			auth.GET("/metrics/series", ListSeries)

			auth.GET("/services", ListMonServices)
			auth.POST("/services", RegisterService)
			auth.PUT("/services/:id", UpdateService)
			auth.DELETE("/services/:id", DeleteService)

			auth.GET("/alerts/rules", ListAlertRules)
			auth.POST("/alerts/rules", CreateAlertRule)
			auth.PUT("/alerts/rules/:id", UpdateAlertRule)
			auth.DELETE("/alerts/rules/:id", DeleteAlertRule)
			auth.GET("/alerts/events", ListAlertEvents)
			auth.POST("/alerts/events/:id/ack", AckAlert)
			auth.POST("/alerts/events/:id/close", CloseAlert)

			auth.GET("/dashboards", ListDashboards)
			auth.POST("/dashboards", CreateDashboard)
			auth.PUT("/dashboards/:id", UpdateDashboard)

			auth.GET("/dashboard", Dashboard)
			auth.GET("/analytics/incidents", IncidentAnalytics)
			auth.GET("/analytics/slo", SLOAnalytics)
		}
	}

	// Graceful shutdown
	addr := ":" + strconv.Itoa(cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: r}
	go func() {
		logger.Info("server listening", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("forced shutdown", "error", err)
	}
	logger.Info("server exited")
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info(fmt.Sprintf("%s %s %d", c.Request.Method, c.Request.URL.Path, c.Writer.Status()),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start).String(),
			"ip", c.ClientIP(),
		)
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func readinessCheck(c *gin.Context) {
	status := svc.HealthCheck()
	code := http.StatusOK
	for _, v := range status {
		if v != "healthy" {
			code = http.StatusServiceUnavailable
			break
		}
	}
	c.JSON(code, gin.H{"checks": status})
}

func Login(c *gin.Context)   { c.JSON(http.StatusOK, gin.H{"message": "login"}) }

func IngestTrace(c *gin.Context)        { c.JSON(http.StatusCreated, gin.H{"message": "trace ingested"}) }
func IngestTraces(c *gin.Context)       { c.JSON(http.StatusCreated, gin.H{"message": "traces ingested"}) }
func ListTraces(c *gin.Context)         { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func GetTrace(c *gin.Context)           { c.JSON(http.StatusOK, gin.H{"data": gin.H{}}) }
func ListTraceSpans(c *gin.Context)     { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func ListTraceServices(c *gin.Context)  { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }

func IngestLog(c *gin.Context)          { c.JSON(http.StatusCreated, gin.H{"message": "log ingested"}) }
func IngestLogs(c *gin.Context)         { c.JSON(http.StatusCreated, gin.H{"message": "logs ingested"}) }
func QueryLogs(c *gin.Context)          { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func ListLogFields(c *gin.Context)      { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func LogAggregations(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }

func IngestMetric(c *gin.Context)       { c.JSON(http.StatusCreated, gin.H{"message": "metric ingested"}) }
func QueryMetrics(c *gin.Context)       { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func ListMetricLabels(c *gin.Context)   { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func ListSeries(c *gin.Context)         { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }

func ListMonServices(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func RegisterService(c *gin.Context)    { c.JSON(http.StatusCreated, gin.H{"message": "service registered"}) }
func UpdateService(c *gin.Context)      { c.JSON(http.StatusOK, gin.H{"message": "service updated"}) }
func DeleteService(c *gin.Context)      { c.JSON(http.StatusOK, gin.H{"message": "service deleted"}) }

func ListAlertRules(c *gin.Context)     { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func CreateAlertRule(c *gin.Context)    { c.JSON(http.StatusCreated, gin.H{"message": "rule created"}) }
func UpdateAlertRule(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"message": "rule updated"}) }
func DeleteAlertRule(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"message": "rule deleted"}) }
func ListAlertEvents(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func AckAlert(c *gin.Context)           { c.JSON(http.StatusOK, gin.H{"message": "ack"}) }
func CloseAlert(c *gin.Context)         { c.JSON(http.StatusOK, gin.H{"message": "closed"}) }

func ListDashboards(c *gin.Context)     { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func CreateDashboard(c *gin.Context)    { c.JSON(http.StatusCreated, gin.H{"message": "dashboard created"}) }
func UpdateDashboard(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"message": "dashboard updated"}) }

func Dashboard(c *gin.Context)          { c.JSON(http.StatusOK, gin.H{"data": gin.H{}}) }
func IncidentAnalytics(c *gin.Context)  { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func SLOAnalytics(c *gin.Context)       { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
