package config

import (
	"encoding/base64"
)

// 这些变量在编译时通过 -ldflags 注入（base64编码）
// GitHub Actions会从secrets读取并注入这些值
var (
	embeddedAgConnectJSON string // base64 encoded
	embeddedPrivateJSON   string // base64 encoded
	embeddedEncryptionKey string // base64 encoded
)

// GetEmbeddedAgConnectJSON 返回嵌入的agconnect配置（解码后）
func GetEmbeddedAgConnectJSON() string {
	if embeddedAgConnectJSON == "" {
		return ""
	}
	decoded, err := base64.StdEncoding.DecodeString(embeddedAgConnectJSON)
	if err != nil {
		return "" // 解码失败返回空
	}
	return string(decoded)
}

// GetEmbeddedPrivateJSON 返回嵌入的private配置（解码后）
func GetEmbeddedPrivateJSON() string {
	if embeddedPrivateJSON == "" {
		return ""
	}
	decoded, err := base64.StdEncoding.DecodeString(embeddedPrivateJSON)
	if err != nil {
		return "" // 解码失败返回空
	}
	return string(decoded)
}

// GetEmbeddedEncryptionKey 返回嵌入的加密密钥
func GetEmbeddedEncryptionKey() string {
	if embeddedEncryptionKey == "" {
		return ""
	}
	decoded, err := base64.StdEncoding.DecodeString(embeddedEncryptionKey)
	if err != nil {
		return "" // 解码失败返回空
	}
	return string(decoded)
}
