package main

import (
	"log"

	"github.com/dengdeng-harmonyos/server/internal/config"
	"github.com/dengdeng-harmonyos/server/internal/database"
	"github.com/dengdeng-harmonyos/server/internal/handler"
	"github.com/dengdeng-harmonyos/server/internal/logger"
	"github.com/dengdeng-harmonyos/server/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志系统
	logger.Init()
	logger.Info("=== Dengdeng Push Server Starting ===")

	// 加载配置（使用嵌入的配置和环境变量）
	cfg := config.Load()
	logger.Info("Configuration loaded:")
	logger.Info("  Server Port: %s", cfg.Server.Port)
	logger.Info("  Server Mode: %s", cfg.Server.Mode)
	logger.Info("  Database: %s:%s/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
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

	pushHandler, err := handler.NewPushHandler(db, deviceHandler, cfg.HuaweiPush, cfg.Server.ServerName)
	if err != nil {
		logger.Error("Failed to create push handler: %v", err)
		log.Fatalf("Failed to create push handler: %v", err)
	}
	logger.Info("✓ Push handler initialized")

	// 创建消息处理器
	messageHandler := handler.NewMessageHandler(db.DB)
	logger.Info("✓ Message handler initialized")

	appUpdateHandler := handler.NewAppUpdateHandler(cfg.AppUpdate)
	logger.Info("✓ App update handler initialized")

	// API v1 路由
	v1 := router.Group("/api/v1")
	{
		// 设备管理
		device := v1.Group("/device")
		{
			device.POST("/register", deviceHandler.Register)       // 注册设备，返回device_id
			device.PUT("/update-token", deviceHandler.UpdateToken) // 更新Push Token
			device.DELETE("/delete", deviceHandler.Delete)         // 删除设备
		}

		// 推送消息（GET方式，方便直接调用）
		push := v1.Group("/push")
		{
			push.GET("/notification", pushHandler.SendNotification) // 发送通知消息
		}

		messages := v1.Group("/messages")
		{
			messages.GET("/pending", messageHandler.GetPendingMessages) // 获取待接收消息
			messages.POST("/confirm", messageHandler.ConfirmMessages)   // 确认消息已收到
		}

		app := v1.Group("/app")
		{
			app.GET("/update", appUpdateHandler.Check) // 检查App强制更新策略
		}
	}

	// 健康检查（支持GET和HEAD）
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"version": "1.0.0",
			"service": "Dengdeng Push Server (Huawei Push Kit v3)",
		})
	})
	router.HEAD("/health", func(c *gin.Context) {
		c.Status(200)
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
