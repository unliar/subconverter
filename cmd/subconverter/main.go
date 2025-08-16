package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	version = "1.0.0"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	var (
		configFile = flag.String("config", "configs/config.yaml", "配置文件路径")
		port       = flag.Int("port", 25500, "服务监听端口")
		showVer    = flag.Bool("version", false, "显示版本信息")
	)
	flag.Parse()

	if *showVer {
		fmt.Printf("SubConverter Go 版本\n")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Date: %s\n", date)
		return
	}

	// 初始化 Gin 路由
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// 基础路由
	router.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": version,
			"commit":  commit,
			"date":    date,
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// 兼容 C++ 版本的核心 API 路由（暂时返回占位符）
	router.GET("/sub", func(c *gin.Context) {
		c.String(http.StatusOK, "# SubConverter Go 版本\n# 订阅转换功能正在开发中...\n")
	})

	router.GET("/getruleset", func(c *gin.Context) {
		c.String(http.StatusOK, "# 规则集功能正在开发中...\n")
	})

	router.GET("/getprofile", func(c *gin.Context) {
		c.String(http.StatusOK, "# 配置档案功能正在开发中...\n")
	})

	// 启动 HTTP 服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: router,
	}

	// 优雅关闭
	go func() {
		log.Printf("SubConverter Go 版本启动，监听端口：%d", *port)
		log.Printf("配置文件：%s", *configFile)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 优雅关闭，5秒超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}

	log.Println("服务器已退出")
}