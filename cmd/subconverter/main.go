package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"subconverter-go/internal/api"
	"subconverter-go/internal/config"
)

var (
	version = "1.0.0"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	var (
		configDir = flag.String("config", "configs", "配置目录路径")
		port      = flag.Int("port", 25500, "服务监听端口")
		showVer   = flag.Bool("version", false, "显示版本信息")
	)
	flag.Parse()

	if *showVer {
		fmt.Printf("SubConverter Go 版本\n")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Date: %s\n", date)
		return
	}

	// 初始化配置管理器
	configManager := config.NewManager(*configDir)
	if err := configManager.LoadConfig(); err != nil {
		log.Printf("Warning: Failed to load config: %v", err)
		log.Printf("Using default configuration")
		// 继续使用现有的configManager，它会使用默认值
	}

	// 创建API服务器
	server, err := api.NewServer(configManager)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// 启动服务器（在goroutine中）
	go func() {
		address := fmt.Sprintf(":%d", *port)
		log.Printf("SubConverter Go 版本启动，监听端口：%d", *port)
		log.Printf("配置目录：%s", *configDir)
		if err := server.Start(address); err != nil {
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
	if err := server.Stop(ctx); err != nil {
		log.Printf("服务器关闭出错: %v", err)
	}

	log.Println("服务器已退出")
}
