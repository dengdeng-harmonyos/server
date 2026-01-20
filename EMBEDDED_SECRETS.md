# 敏感配置嵌入说明

## 实现方式

已将 `agconnect-services.json` 和 `private.json` 的内容嵌入到 Go 源代码中：

### 文件位置
- **嵌入配置文件**: `server/internal/config/embedded_secrets.go`
- **配置加载逻辑**: `server/internal/config/config.go`
- **推送服务加载**: `server/internal/service/huawei_push.go`

### 工作原理

1. **编译时嵌入**
   - 敏感配置以字符串常量形式存储在 `embedded_secrets.go` 中
   - 编译后的二进制文件中，这些字符串会被编码，不是明文存储

2. **运行时加载优先级**
   ```
   嵌入配置（优先） → 文件配置（降级）
   ```
   - 优先使用嵌入在二进制中的配置
   - 如果嵌入配置为空，则从文件读取（用于开发环境）

3. **安全性**
   - 编译后的二进制文件中，字符串被混淆存储
   - 无法通过简单的文本查看工具读取原始JSON
   - 需要使用strings命令或反编译工具才能提取，增加了攻击难度

## 使用方法

### 开发环境
保留配置文件用于开发：
```bash
config/
├── agconnect-services.json
└── private.json
```

### 生产环境
1. 编译二进制：
   ```bash
   go build -o bin/server cmd/server/main.go
   ```

2. 可以删除配置文件（可选）：
   ```bash
   rm config/agconnect-services.json
   rm config/private.json
   ```

3. 运行服务：
   ```bash
   ./bin/server
   ```
   服务会自动使用嵌入的配置。

## 更新配置

如果需要更新配置（如更换密钥），修改 `embedded_secrets.go` 中的常量：

```go
const (
    embeddedAgConnectJSON = `{新的配置内容}`
    embeddedPrivateJSON = `{新的配置内容}`
)
```

然后重新编译。

## 安全建议

1. **版本控制**
   - 将 `embedded_secrets.go` 添加到 `.gitignore`
   - 或者使用 Git Crypt 等工具加密该文件

2. **进一步混淆**（可选）
   - 使用 Base64 编码配置字符串
   - 使用 XOR 等简单加密算法
   - 在运行时解密

3. **环境变量优先**
   - 生产环境建议使用环境变量覆盖配置
   - 这样可以避免在代码中硬编码敏感信息

## 示例：添加混淆

如果需要更高的安全性，可以对配置进行Base64编码：

```go
// embedded_secrets.go
import "encoding/base64"

const (
    embeddedAgConnectJSONBase64 = "eyJhZ2N...编码后的内容..."
    embeddedPrivateJSONBase64 = "eyJwcm9...编码后的内容..."
)

func GetEmbeddedAgConnectJSON() string {
    decoded, _ := base64.StdEncoding.DecodeString(embeddedAgConnectJSONBase64)
    return string(decoded)
}
```

这样在二进制文件中就不会出现明文JSON结构。
