package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yourusername/dangdangdang-push-server/internal/config"
	"github.com/yourusername/dangdangdang-push-server/internal/database"
	"github.com/yourusername/dangdangdang-push-server/internal/handler"
	"github.com/yourusername/dangdangdang-push-server/internal/logger"
	"github.com/yourusername/dangdangdang-push-server/internal/middleware"
)

func main() {
	// 初始化日志系统
	logger.Init()
	logger.Info("=== Dangdangdang Push Server Starting ===")

	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using environment variables")
	} else {
		logger.Info("Loaded .env file successfully")
	}

	// 加载配置
	cfg := config.Load()
	logger.Info("Configuration loaded:")
	logger.Info("  Server Port: %s", cfg.Server.Port)
	logger.Info("  Server Mode: %s", cfg.Server.Mode)
	logger.Info("  Database: %s:%s/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	logger.Info("  Huawei Project ID: %s", cfg.HuaweiPush.ProjectID)
	logger.Info("  Huawei Push API: %s", cfg.HuaweiPush.PushAPIURL)

	// 初始化数据库
	logger.Info("Connecting to database...")
	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database: %v", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	logger.Info("✓ Database connected successfully")

	// 初始化数据库表
	logger.Info("Initializing database tables...")
	if err := db.InitTables(); err != nil {
		logger.Error("Failed to initialize tables: %v", err)
		log.Fatalf("Failed to initialize tables: %v", err)
	}
	logger.Info("✓ Database tables initialized")

	// 设置 Gin 模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由（不使用默认中间件）
	router := gin.New()

	// 使用自定义中间件
	router.Use(logger.GinRecovery())
	router.Use(logger.GinLogger())
	router.Use(middleware.CORS())

	// 初始化处理器
	logger.Info("Initializing handlers...")
	deviceHandler, err := handler.NewDeviceHandler(db, *cfg)
	if err != nil {
		logger.Error("Failed to create device handler: %v", err)
		log.Fatalf("Failed to create device handler: %v", err)
	}
	logger.Info("✓ Device handler initialized")

	pushHandler, err := handler.NewPushHandler(db, deviceHandler, cfg.HuaweiPush)
	if err != nil {
		logger.Error("Failed to create push handler: %v", err)
		log.Fatalf("Failed to create push handler: %v", err)
	}
	logger.Info("✓ Push handler initialized")

	// API v1 路由
	v1 := router.Group("/api/v1")
	{
		// 设备管理
		device := v1.Group("/device")
		{
			device.POST("/register", deviceHandler.Register)       // 注册设备，返回device_key
			device.PUT("/update-token", deviceHandler.UpdateToken) // 更新Push Token
			device.GET("/deactivate", deviceHandler.Deactivate)    // 停用设备
		}

		// 推送消息（GET方式，方便直接调用）
		push := v1.Group("/push")
		{
			push.GET("/notification", pushHandler.SendNotification)    // 发送通知消息
			push.GET("/form", pushHandler.SendFormUpdate)              // 发送卡片刷新消息
			push.GET("/background", pushHandler.SendBackgroundMessage) // 发送后台消息
			push.GET("/batch", pushHandler.SendBatch)                  // 批量推送
			push.GET("/statistics", pushHandler.GetStatistics)         // 查询统计数据
		}
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"version": "1.0.0",
			"service": "Dangdangdang Push Server (Huawei Push Kit v3)",
		})
	})

	// 启动服务器
	logger.Info("===========================================")
	logger.Info("🚀 Server is ready!")
	logger.Info("   Listening on: http://0.0.0.0:%s", cfg.Server.Port)
	logger.Info("   Health check: http://0.0.0.0:%s/health", cfg.Server.Port)
	logger.Info("   Push endpoint: http://0.0.0.0:%s/api/v1/push/notification", cfg.Server.Port)
	logger.Info("===========================================")

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		logger.Error("Failed to start server: %v", err)
		log.Fatalf("Failed to start server: %v", err)
	}
}
