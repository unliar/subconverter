package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"subconverter-go/internal/config"
	"subconverter-go/internal/service"
)

// Server HTTP服务器
type Server struct {
	Router       *gin.Engine // 导出router字段用于测试
	httpServer   *http.Server
	config       *config.Manager
	converterSvc *service.ConverterService
	rulesetSvc   *service.RulesetService
	logger       *logrus.Logger
}

// NewServer 创建新的服务器实例
func NewServer(configManager *config.Manager) (*Server, error) {
	// 初始化日志
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	
	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)
	
	// 创建路由器
	router := gin.New()
	
	// 创建服务层
	converterSvc := service.NewConverterService(configManager, logger)
	rulesetSvc := service.NewRulesetService(configManager, logger)
	
	server := &Server{
		Router:       router,
		config:       configManager,
		converterSvc: converterSvc,
		rulesetSvc:   rulesetSvc,
		logger:       logger,
	}
	
	// 设置路由
	server.setupRoutes()
	
	return server, nil
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 添加中间件
	s.Router.Use(s.recoveryMiddleware())
	s.Router.Use(s.corsMiddleware())
	s.Router.Use(s.loggerMiddleware())
	
	// 健康检查
	s.Router.GET("/health", s.handleHealth)
	s.Router.GET("/version", s.handleVersion)
	
	// 核心API - 兼容 C++ 版本
	s.Router.GET("/sub", s.handleSubscriptionConvert)
	s.Router.POST("/sub", s.handleSubscriptionConvert)
	
	// 规则集API
	s.Router.GET("/getruleset", s.handleGetRuleset)
	s.Router.GET("/refreshrules", s.handleRefreshRules)
	s.Router.GET("/rulesets", s.handleGetAvailableRulesets)
	
	// 配置API
	s.Router.GET("/getprofile", s.handleGetProfile)
	s.Router.POST("/updateconf", s.handleUpdateConfig)
	s.Router.GET("/readconf", s.handleReadConfig)
	
	// 短链接API
	s.Router.GET("/short", s.handleShortURL)
	s.Router.POST("/short", s.handleCreateShortURL)
	
	// 工具API
	s.Router.GET("/targets", s.handleGetSupportedTargets)
	
	// 静态文件服务
	s.Router.Static("/assets", "./assets")
	s.Router.StaticFile("/", "./assets/index.html")
}

// Start 启动服务器
func (s *Server) Start(address string) error {
	// 创建 HTTP 服务器
	s.httpServer = &http.Server{
		Addr:           address,
		Handler:        s.Router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
	
	s.logger.Infof("Starting server on %s", address)
	
	// 启动服务器
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}
	
	return nil
}

// Stop 停止服务器
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Shutting down server...")
	
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	
	return nil
}

// recoveryMiddleware 恢复中间件
func (s *Server) recoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			s.logger.Errorf("Panic recovered: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// corsMiddleware CORS中间件
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// loggerMiddleware 日志中间件
func (s *Server) loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		
		// 处理请求
		c.Next()
		
		// 记录日志
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		
		if raw != "" {
			path = path + "?" + raw
		}
		
		s.logger.WithFields(logrus.Fields{
			"status":   statusCode,
			"latency":  latency,
			"clientIP": clientIP,
			"method":   method,
			"path":     path,
		}).Info("Request processed")
	}
}

// handleHealth 处理健康检查
func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().Unix(),
	})
}

// handleVersion 处理版本信息
func (s *Server) handleVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":        "SubConverter-Go",
		"version":     "1.0.0",
		"description": "Go version of SubConverter",
	})
}
