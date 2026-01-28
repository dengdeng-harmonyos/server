package models

import "time"

// Device 设备信息（简化版）
type Device struct {
	ID           int       `json:"id"`
	DeviceId     string    `json:"device_id"`  // 服务端生成的随机ID（对外使用）
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
	Success    bool   `json:"success"`
	DeviceId   string `json:"device_id"`
	ServerName string `json:"server_name"` // 服务器名称
	Message    string `json:"message"`
}

// PushNotificationRequest 通知消息推送请求（GET参数）
type PushNotificationRequest struct {
	DeviceId  string `form:"device_id" binding:"required"`
	Title     string `form:"title" binding:"required"`
	Content   string `form:"content" binding:"required"`
	Data      string `form:"data"` // JSON字符串
}

// FormUpdateRequest 卡片刷新请求（GET参数）
type FormUpdateRequest struct {
	DeviceId  string `form:"device_id" binding:"required"`
	FormID    string `form:"form_id" binding:"required"`
	FormData  string `form:"form_data" binding:"required"` // JSON字符串
}

// BackgroundPushRequest 后台推送请求（GET参数）
type BackgroundPushRequest struct {
	DeviceId  string `form:"device_id" binding:"required"`
	Data      string `form:"data" binding:"required"` // JSON字符串
}

// BatchPushRequest 批量推送请求（GET参数）
type BatchPushRequest struct {
	DeviceIds  string `form:"device_ids" binding:"required"` // 逗号分隔的device_id列表
	Title      string `form:"title" binding:"required"`
	Body       string `form:"body" binding:"required"`
	Data       string `form:"data"` // JSON字符串
}

// UnifiedApiResponse 统一API响应格式
type UnifiedApiResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// ErrorCode 错误码定义
type ErrorCode struct{}

var ErrorCodes = ErrorCode{}

// 系统错误 1xxx
const (
	SystemError   = 1000 // 系统错误
	InvalidParams = 1001 // 参数无效
)

// 认证错误 2xxx
const (
	Unauthorized       = 2001 // 未授权
	InvalidSignature   = 2002 // 签名无效
	SignatureExpired   = 2003 // 签名过期
	InvalidAppID       = 2004 // 无效的AppID
	VersionTooOld      = 2005 // 版本过旧
)

// 业务错误 3xxx
const (
	BusinessError     = 3001 // 业务错误
	ResourceNotFound  = 3002 // 资源未找到
	OperationFailed   = 3003 // 操作失败
)

// 数据错误 4xxx
const (
	DataNotFound      = 4001 // 数据未找到
	DataAlreadyExists = 4002 // 数据已存在
)
