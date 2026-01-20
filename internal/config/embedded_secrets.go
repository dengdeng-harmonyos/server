package config

// 这些变量在编译时通过 -ldflags 注入
// GitHub Actions会从secrets读取并注入这些值
var (
	embeddedAgConnectJSON string
	embeddedPrivateJSON   string
	embeddedEncryptionKey string
)

// GetEmbeddedAgConnectJSON 返回嵌入的agconnect配置
func GetEmbeddedAgConnectJSON() string {
	return embeddedAgConnectJSON
}

// GetEmbeddedPrivateJSON 返回嵌入的private配置
func GetEmbeddedPrivateJSON() string {
	return embeddedPrivateJSON
}

// GetEmbeddedEncryptionKey 返回嵌入的加密密钥
func GetEmbeddedEncryptionKey() string {
	return embeddedEncryptionKey
}
