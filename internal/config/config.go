package config

import (
	"os"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	HuaweiPush HuaweiPushConfig
	Security   SecurityConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type HuaweiPushConfig struct {
	ProjectID          string
	ServiceAccountFile string // JWT私钥文件路径
	JWTExpiry          int    // JWT过期时间（秒）
	PushAPIURL         string
}

type SecurityConfig struct {
	EncryptionKey         string // Push Token加密密钥（32字节）
	DeviceKeyTTL          int    // Device Key有效期（秒）
	MaxDailyPushPerDevice int    // 每设备每日最大推送数
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "push_server"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		HuaweiPush: HuaweiPushConfig{
			ProjectID:          getEnv("HUAWEI_PROJECT_ID", ""),
			ServiceAccountFile: getEnv("HUAWEI_SERVICE_ACCOUNT_FILE", "./config/agconnect-services.json"),
			JWTExpiry:          3600,
			PushAPIURL:         "https://push-api.cloud.huawei.com/v3",
		},
		Security: SecurityConfig{
			EncryptionKey:         getEnv("PUSH_TOKEN_ENCRYPTION_KEY", ""),
			DeviceKeyTTL:          2592000, // 30天
			MaxDailyPushPerDevice: 100,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
