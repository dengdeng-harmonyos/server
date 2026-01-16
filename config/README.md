# 配置文件说明

## agconnect-services.json

这是华为AppGallery Connect的配置文件，包含项目的认证信息。

### 如何获取

1. 登录 [AppGallery Connect](https://developer.huawei.com/consumer/cn/service/josp/agc/index.html)
2. 选择您的项目
3. 进入 **项目设置 → 常规 → 我的应用**
4. 点击 **下载 agconnect-services.json**
5. 将下载的文件保存到此目录

### 文件结构

```json
{
  "client": {
    "project_id": "101653523863440882",  // 项目ID - 需要配置到.env
    "app_id": "6917594697233331968",      // 应用ID
    "client_id": "...",                    // 客户端ID
    "client_secret": "...",                // 客户端密钥
    "api_key": "...",                      // API密钥
    "package_name": "your.package.name"    // 应用包名
  },
  "service": {
    // 华为各项服务的配置
  }
}
```

### 重要字段

- **client.project_id**: 项目ID，需要配置到 `.env` 的 `HUAWEI_PROJECT_ID`
- **client.app_id**: 应用ID
- **client.package_name**: 应用包名，需与HarmonyOS应用配置一致

### 安全提醒

⚠️ **重要**: 此文件包含敏感信息，请勿提交到Git仓库！

`.gitignore` 已配置忽略此文件：
```
config/agconnect-services.json
```

## 相关文档

- [华为Push Kit文档](https://developer.huawei.com/consumer/cn/doc/HMSCore-Guides/service-introduction-0000001050040060)
- [AppGallery Connect配置](https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-Guides/agc-get-started-harmonyos-0000001933963166)
