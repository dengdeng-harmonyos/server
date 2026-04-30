package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	HuaweiPush HuaweiPushConfig
	Security   SecurityConfig
	AppUpdate  AppUpdateConfig
}

type ServerConfig struct {
	Port         string
	Mode         string
	ServerName   string // 服务器名称，用于标识消息来源
	Version      string
	Build        string
	APIVersion   int64
	Capabilities []string
	UpgradeURL   string
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
	DeviceIdTTL           int    // Device Id有效期（秒）
	MaxDailyPushPerDevice int    // 每设备每日最大推送数
}

type AppUpdateConfig struct {
	LatestVersionCode int64
	LatestVersionName string
	MinVersionCode    int64
	ForceUpdate       bool
	StoreURL          string
	ReleaseNotes      string
	PolicyFile        string
}

// AgConnectServices 用于解析agconnect-services.json
type AgConnectServices struct {
	Client struct {
		ProjectID string `json:"project_id"`
	} `json:"client"`
}

// loadProjectIDFromAgConnect 从嵌入的配置读取项目ID
func loadProjectIDFromAgConnect() (string, error) {
	// 使用嵌入的配置（编译时注入）
	embeddedJSON := GetEmbeddedAgConnectJSON()
	if embeddedJSON == "" {
		return "", fmt.Errorf("embedded agconnect configuration is empty")
	}

	var agConnect AgConnectServices
	if err := json.Unmarshal([]byte(embeddedJSON), &agConnect); err != nil {
		return "", fmt.Errorf("failed to parse embedded agconnect configuration: %w", err)
	}

	if agConnect.Client.ProjectID == "" {
		return "", fmt.Errorf("project_id not found in embedded configuration")
	}

	return agConnect.Client.ProjectID, nil
}

func Load() *Config {
	// 尝试从嵌入配置读取ProjectID
	projectID := getEnv("HUAWEI_PROJECT_ID", "")

	// 如果环境变量未设置，从嵌入配置读取
	if projectID == "" {
		if pid, err := loadProjectIDFromAgConnect(); err == nil {
			projectID = pid
		}
	}

	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			Mode:         getEnv("GIN_MODE", "debug"),
			ServerName:   getEnv("SERVER_NAME", "噔噔推送服务"),
			Version:      getEnv("SERVER_VERSION", "1.1.0"),
			Build:        getEnv("SERVER_BUILD", "dev"),
			APIVersion:   getEnvInt64("SERVER_API_VERSION", 2),
			Capabilities: getEnvStringList("SERVER_CAPABILITIES", []string{"message_crypto_v1", "push_url_data", "app_update_policy"}),
			UpgradeURL:   getEnv("SERVER_UPGRADE_URL", "https://github.com/dengdeng-harmonyos/server"),
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
			ProjectID:          projectID,
			ServiceAccountFile: "", // 不再使用文件，配置已嵌入
			JWTExpiry:          3600,
			PushAPIURL:         "https://push-api.cloud.huawei.com/v3",
		},
		Security: SecurityConfig{
			EncryptionKey:         getEncryptionKey(),
			DeviceIdTTL:           2592000, // 30天
			MaxDailyPushPerDevice: 100,
		},
		AppUpdate: AppUpdateConfig{
			LatestVersionCode: getEnvInt64("APP_LATEST_VERSION_CODE", 0),
			LatestVersionName: getEnv("APP_LATEST_VERSION_NAME", ""),
			MinVersionCode:    getEnvInt64("APP_MIN_VERSION_CODE", 0),
			ForceUpdate:       getEnvBool("APP_FORCE_UPDATE", true),
			StoreURL:          getEnv("APP_STORE_URL", "store://appgallery.huawei.com/app/detail?id=top.yidingyaojizhu.dengdeng"),
			ReleaseNotes:      getEnv("APP_RELEASE_NOTES", ""),
			PolicyFile:        getEnv("APP_UPDATE_POLICY_FILE", "config/app_update_policy.json"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	var parsed int64
	if _, err := fmt.Sscanf(value, "%d", &parsed); err != nil {
		return defaultValue
	}
	return parsed
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	switch value {
	case "1", "true", "TRUE", "True", "yes", "YES", "Yes", "on", "ON", "On":
		return true
	case "0", "false", "FALSE", "False", "no", "NO", "No", "off", "OFF", "Off":
		return false
	default:
		return defaultValue
	}
}

func getEnvStringList(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	var result []string
	var current string
	for _, char := range value {
		if char == ',' {
			if current != "" {
				result = append(result, current)
			}
			current = ""
			continue
		}
		if char != ' ' && char != '\t' && char != '\n' && char != '\r' {
			current += string(char)
		}
	}

	if current != "" {
		result = append(result, current)
	}
	if len(result) == 0 {
		return defaultValue
	}
	return result
}

// getEncryptionKey 获取加密密钥（优先使用环境变量）
func getEncryptionKey() string {
	// 优先使用环境变量（运行时配置）
	envKey := os.Getenv("PUSH_TOKEN_ENCRYPTION_KEY")
	if envKey != "" {
		return envKey
	}

	// 降级到嵌入的密钥（编译时注入）
	embedded := GetEmbeddedEncryptionKey()
	if embedded != "" {
		return embedded
	}

	// 最后返回空字符串
	return ""
}
