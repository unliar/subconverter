// internal/api/server.go
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
	"subconverter-go/pkg/constants"
)

// Server HTTP服务器
type Server struct {
	router         *gin.Engine
	httpServer     *http.Server
	config         *config.Manager
	converterSvc   *service.ConverterService
	rulesetSvc     *service.RulesetService
	logger         *logrus.Logger
}

// NewServer 创建新的服务器实例
func NewServer(configManager *config.Manager) (*Server, error) {
	// 初始化日志
	logger := logrus.New()
	
	// 设置 Gin 模式
	appConfig := configManager.GetApp()
	if appConfig != nil && appConfig.Log != nil {
		switch appConfig.Log.Level {
		case "debug":
			gin.SetMode(gin.DebugMode)
			logger.SetLevel(logrus.DebugLevel)
		case "error", "warn":
			gin.SetMode(gin.ReleaseMode)
			logger.SetLevel(logrus.WarnLevel)
		default:
			gin.SetMode(gin.ReleaseMode)
			logger.SetLevel(logrus.InfoLevel)
		}
	}
	
	// 创建路由器
	router := gin.New()
	
	// 创建服务层
	converterSvc := service.NewConverterService(configManager, logger)
	rulesetSvc := service.NewRulesetService(configManager, logger)
	
	server := &Server{
		router:       router,
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
	s.router.Use(s.recoveryMiddleware())
	s.router.Use(s.corsMiddleware())
	s.router.Use(s.loggerMiddleware())
	
	// 检查是否需要认证
	appConfig := s.config.GetApp()
	if appConfig != nil && appConfig.Server != nil && appConfig.Server.AccessToken != "" {
		s.router.Use(s.authMiddleware())
	}
	
	// 健康检查
	s.router.GET("/health", s.handleHealth)
	s.router.GET("/version", s.handleVersion)
	
	// 核心API - 完全兼容 C++ 版本
	s.router.GET("/sub", s.handleSubscriptionConvert)
	s.router.POST("/sub", s.handleSubscriptionConvert)
	
	// 规则集API
	s.router.GET("/getruleset", s.handleGetRuleset)
	s.router.GET("/refreshrules", s.handleRefreshRules)
	
	// 配置API
	s.router.GET("/getprofile", s.handleGetProfile)
	s.router.POST("/updateconf", s.handleUpdateConfig)
	s.router.GET("/readconf", s.handleReadConfig)
	
	// 短链接API
	s.router.GET("/short", s.handleShortURL)
	s.router.POST("/short", s.handleCreateShortURL)
	
	// 静态文件服务
	s.router.Static("/assets", "./assets")
	s.router.StaticFile("/", "./assets/index.html")
}

// Start 启动服务器
func (s *Server) Start(address string) error {
	appConfig := s.config.GetApp()
	
	// 创建 HTTP 服务器
	s.httpServer = &http.Server{
		Addr:    address,
		Handler: s.router,
	}
	
	// 设置超时
	if appConfig != nil && appConfig.Server != nil {
		s.httpServer.ReadTimeout = appConfig.Server.ReadTimeout
		s.httpServer.WriteTimeout = appConfig.Server.WriteTimeout
		s.httpServer.MaxHeaderBytes = 1 << 20 // 1 MB
	} else {
		s.httpServer.ReadTimeout = 30 * time.Second
		s.httpServer.WriteTimeout = 30 * time.Second
		s.httpServer.MaxHeaderBytes = 1 << 20
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

// authMiddleware 认证中间件
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appConfig := s.config.GetApp()
		if appConfig == nil || appConfig.Server == nil || appConfig.Server.AccessToken == "" {
			c.Next()
			return
		}
		
		// 获取访问令牌
		token := c.Query("access_token")
		if token == "" {
			token = c.GetHeader("Authorization")
			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}
		}
		
		// 验证令牌
		if token != appConfig.Server.AccessToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		
		c.Next()
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
		"name":        constants.AppName,
		"version":     "1.0.0",
		"description": constants.AppDescription,
	})
}