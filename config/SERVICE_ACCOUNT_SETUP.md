# 华为Push Kit服务账号配置说明

## ⚠️ 重要：需要创建服务账号密钥文件

华为Push Kit使用JWT token认证，需要从华为开发者联盟下载服务账号密钥文件。

## 📝 创建步骤

### 1. 访问API Console
登录华为开发者联盟，访问：https://developer.huawei.com/consumer/cn/console/api/myApi

### 2. 选择项目
选择你的应用所属的项目（Project ID: 101653523863440882）

### 3. 创建服务账号密钥
- 点击"创建凭证"
- 选择"服务账号密钥"
- 下载JSON文件

### 4. 保存密钥文件
将下载的JSON文件保存为：`config/private.json`

## 📄 密钥文件格式示例

```json
{
  "project_id": "101653523863440882",
  "key_id": "xxxxxxxxxx",
  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIJQgIBADANBgkqhkiG9w0...\n-----END PRIVATE KEY-----\n",
  "sub_account": "xxxxxxxxxx",
  "auth_uri": "https://oauth-login.cloud.huawei.com/oauth2/v3/authorize",
  "token_uri": "https://oauth-login.cloud.huawei.com/oauth2/v3/token",
  "auth_provider_cert_uri": "https://oauth-login.cloud.huawei.com/oauth2/v3/certs",
  "client_cert_uri": "https://oauth-login.cloud.huawei.com/oauth2/v3/x509?client_id="
}
```

## ⚙️ 更新配置

### Docker环境变量
在 `docker-compose.single.yml` 中设置：

```yaml
environment:
  - HUAWEI_SERVICE_ACCOUNT_FILE=/app/config/private.json
```

### 本地开发
直接将文件保存到：`./config/private.json`

## 🔍 验证

文件创建后，重新构建并启动服务：

```bash
./rebuild.sh
```

查看日志应该显示：
```
[INFO] Initializing Huawei Push Service...
[DEBUG] Loading service account from: ./config/private.json
[INFO] ✓ Huawei Push service account loaded
[DEBUG] Key ID: xxxxxxxxxx
[DEBUG] Sub Account: xxxxxxxxxx
[DEBUG] Project ID: 101653523863440882
```

## 📚 参考文档

- [华为Push Kit JWT Token文档](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/push-jwt-token)
- [API服务操作指南](https://developer.huawei.com/consumer/cn/doc/start/api-0000001062522591)

## ⚠️ 注意事项

1. **私钥安全**：服务账号密钥包含私钥，请妥善保管，不要提交到版本控制系统
2. **项目ID匹配**：确保密钥文件中的project_id与你的应用所属项目一致
3. **有效期**：JWT token有效期为3600秒（1小时），系统会自动刷新
4. **时间同步**：服务器时间需要校准为标准时间

## 🔐 安全建议

在生产环境中，建议：
- 使用环境变量存储敏感信息
- 定期轮换服务账号密钥
- 限制密钥文件的访问权限（chmod 600）
