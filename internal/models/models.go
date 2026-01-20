package models

import "time"

// Device 设备信息（简化版）
type Device struct {
	ID           int       `json:"id"`
	DeviceKey    string    `json:"device_key"`  // 服务端生成的随机ID（对外使用）
	PushToken    string    `json:"-"`           // 华为Push Token（加密存储，不对外暴露）
	PublicKey    string    `json:"-"`           // RSA公钥(PEM格式，不对外暴露)
	DeviceType   string    `json:"device_type"` // phone/tablet/watch
	OSVersion    string    `json:"os_version"`  // HarmonyOS版本
	AppVersion   string    `json:"app_version"` // App版本
	IsActive     bool      `json:"is_active"`
	LastActiveAt time.Time `json:"last_active_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PushStatistics 推送统计（仅统计数据，不记录内容）
type PushStatistics struct {
	ID           int       `json:"id"`
	Date         time.Time `json:"date"`
	PushType     string    `json:"push_type"` // notification/form/background
	TotalCount   int       `json:"total_count"`
	SuccessCount int       `json:"success_count"`
	FailedCount  int       `json:"failed_count"`
	CreatedAt    time.Time `json:"created_at"`
}

// DeviceRegisterRequest 设备注册请求
type DeviceRegisterRequest struct {
	PushToken  string `json:"push_token" binding:"required"`
	PublicKey  string `json:"public_key"` // RSA公钥(PEM格式)
	DeviceType string `json:"device_type"`
	OSVersion  string `json:"os_version"`
	AppVersion string `json:"app_version"`
}

// DeviceRegisterResponse 设备注册响应
type DeviceRegisterResponse struct {
	Success   bool   `json:"success"`
	DeviceKey string `json:"device_key"`
	Message   string `json:"message"`
}

// PushNotificationRequest 通知消息推送请求（GET参数）
type PushNotificationRequest struct {
	DeviceKey string `form:"device_key" binding:"required"`
	Title     string `form:"title" binding:"required"`
	Body      string `form:"body" binding:"required"`
	Data      string `form:"data"` // JSON字符串
}

// FormUpdateRequest 卡片刷新请求（GET参数）
type FormUpdateRequest struct {
	DeviceKey string `form:"device_key" binding:"required"`
	FormID    string `form:"form_id" binding:"required"`
	FormData  string `form:"form_data" binding:"required"` // JSON字符串
}

// BackgroundPushRequest 后台推送请求（GET参数）
type BackgroundPushRequest struct {
	DeviceKey string `form:"device_key" binding:"required"`
	Data      string `form:"data" binding:"required"` // JSON字符串
}

// BatchPushRequest 批量推送请求（GET参数）
type BatchPushRequest struct {
	DeviceKeys string `form:"device_keys" binding:"required"` // 逗号分隔的device_key列表
	Title      string `form:"title" binding:"required"`
	Body       string `form:"body" binding:"required"`
	Data       string `form:"data"` // JSON字符串
}
