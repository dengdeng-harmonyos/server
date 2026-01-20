# GitHub Actions 配置指南

本项目使用GitHub Actions自动构建，并通过Secrets安全地注入敏感配置。

## 设置步骤

### 1. 配置GitHub Secrets

进入你的GitHub仓库 → Settings → Secrets and variables → Actions，添加以下secrets：

#### 必需的Secrets

**AGCONNECT_JSON**
```json
{"agcgw_all":{"SG":"connect-dra.dbankcloud.cn","SG_back":"connect-dra.hispace.hicloud.com","CN":"connect-drcn.dbankcloud.cn","CN_back":"connect-drcn.hispace.dbankcloud.com","RU":"connect-drru.hispace.dbankcloud.ru","RU_back":"connect-drru.hispace.dbankcloud.cn","DE":"connect-dre.dbankcloud.cn","DE_back":"connect-dre.hispace.hicloud.com"},"websocketgw_all":{"SG":"connect-ws-dra.hispace.dbankcloud.cn","SG_back":"connect-ws-dra.hispace.dbankcloud.com","CN":"connect-ws-drcn.hispace.dbankcloud.cn","CN_back":"connect-ws-drcn.hispace.dbankcloud.com","RU":"connect-ws-drru.hispace.dbankcloud.ru","RU_back":"connect-ws-drru.hispace.dbankcloud.cn","DE":"connect-ws-dre.hispace.dbankcloud.cn","DE_back":"connect-ws-dre.hispace.dbankcloud.com"},"client":{"cp_id":"YOUR_CP_ID","product_id":"YOUR_PRODUCT_ID","client_id":"YOUR_CLIENT_ID","client_secret":"YOUR_CLIENT_SECRET","project_id":"YOUR_PROJECT_ID","app_id":"YOUR_APP_ID","api_key":"YOUR_API_KEY","package_name":"YOUR_PACKAGE_NAME"},"oauth_client":{"client_id":"YOUR_CLIENT_ID","client_type":30},"app_info":{"app_id":"YOUR_APP_ID","package_name":"YOUR_PACKAGE_NAME"},"code":{"code1":"YOUR_CODE1","code2":"YOUR_CODE2","code3":"YOUR_CODE3","code4":"YOUR_CODE4"},"configuration_version":"3.0","appInfos":[{"package_name":"YOUR_PACKAGE_NAME","client":{"client_secret":"YOUR_CLIENT_SECRET","app_id":"YOUR_APP_ID","api_key":"YOUR_API_KEY"},"code":{"code1":"YOUR_CODE1","code2":"YOUR_CODE2","code3":"YOUR_CODE3","code4":"YOUR_CODE4"}}]}
```
> ⚠️ 将上面的内容复制并替换为你的实际 `agconnect-services.json` 内容（压缩成一行）

**PRIVATE_JSON**
```json
{"project_id":"YOUR_PROJECT_ID","key_id":"YOUR_KEY_ID","private_key":"-----BEGIN PRIVATE KEY-----\nYOUR_PRIVATE_KEY_CONTENT\n-----END PRIVATE KEY-----\n","sub_account":"YOUR_SUB_ACCOUNT","auth_uri":"https://oauth-login.cloud.huawei.com/oauth2/v3/authorize","token_uri":"https://oauth-login.cloud.huawei.com/oauth2/v3/token","auth_provider_cert_uri":"https://oauth-login.cloud.huawei.com/oauth2/v3/certs","client_cert_uri":"https://oauth-login.cloud.huawei.com/oauth2/v3/x509?client_id="}
```
> ⚠️ 将上面的内容替换为你的实际 `private.json` 内容（压缩成一行）

#### 可选的Secrets（用于Docker推送）

- **DOCKER_USERNAME**: 你的Docker Hub用户名
- **DOCKER_PASSWORD**: 你的Docker Hub密码或访问令牌

### 2. 准备Secrets内容

#### 方法1：手动复制粘贴

```bash
# 压缩JSON为一行（去除空格和换行）
cat config/agconnect-services.json | jq -c .
cat config/private.json | jq -c .
```

复制输出的内容到GitHub Secrets。

#### 方法2：使用脚本生成

```bash
# 生成压缩的JSON
cd server
echo "AGCONNECT_JSON:"
cat config/agconnect-services.json | jq -c .
echo ""
echo "PRIVATE_JSON:"
cat config/private.json | jq -c .
```

### 3. 触发构建

#### 自动触发
- 推送代码到 `main` 或 `master` 分支
- 创建Pull Request
- 创建tag（格式：`v*`，如 `v1.0.0`）

#### 手动触发
进入仓库的 Actions 标签页，选择 "Build and Release" 工作流，点击 "Run workflow"。

### 4. 下载构建产物

构建完成后，可以在两个地方获取二进制文件：

1. **Actions Artifacts**（保留30天）
   - 进入 Actions → 选择运行 → Artifacts 部分下载

2. **Releases**（仅限tag触发）
   - 如果推送tag（如 `v1.0.0`），会自动创建Release并附带二进制文件

## 本地构建（开发环境）

如果你有配置文件，可以在本地构建：

```bash
cd server

# 方法1：使用构建脚本（推荐）
chmod +x scripts/build-with-secrets.sh
./scripts/build-with-secrets.sh

# 方法2：手动编译（需要jq工具）
AGCONNECT_JSON=$(cat config/agconnect-services.json | jq -c . | sed 's/"/\\"/g')
PRIVATE_JSON=$(cat config/private.json | jq -c . | sed 's/"/\\"/g')

go build -ldflags "\
  -X 'github.com/dengdeng-harmenyos/server/internal/config.embeddedAgConnectJSON=$AGCONNECT_JSON' \
  -X 'github.com/dengdeng-harmenyos/server/internal/config.embeddedPrivateJSON=$PRIVATE_JSON' \
  -s -w" \
  -o bin/dengdeng-server \
  cmd/server/main.go
```

## 工作流说明

### build.yml 工作流

```yaml
jobs:
  build:    # 构建二进制文件
  docker:   # 构建Docker镜像（可选）
```

**触发条件**：
- `build` job: 所有push、PR、tag都会触发
- `docker` job: 仅在push到main/master分支或tag时触发

**构建选项**：
- `-ldflags -X`: 在编译时注入变量
- `-s -w`: 去除调试信息，减小二进制文件大小

## 安全性说明

✅ **优点**：
1. 敏感信息存储在GitHub Secrets中，加密保存
2. 源代码不包含任何敏感信息，可以公开
3. 编译后的二进制文件包含配置，无需外部文件
4. 只有仓库管理员可以访问Secrets

⚠️ **注意事项**：
1. 不要在日志中打印Secrets内容
2. 二进制文件仍可能被反编译，建议结合其他安全措施
3. 定期轮换密钥和凭证
4. 限制GitHub Actions的权限范围

## 故障排除

### 构建失败

检查以下几点：
1. Secrets是否正确配置
2. JSON格式是否有效（可用 jq 验证）
3. 检查Actions日志中的错误信息

### 测试Secrets格式

```bash
# 验证JSON格式
echo "$AGCONNECT_JSON" | jq .
echo "$PRIVATE_JSON" | jq .

# 如果JSON无效，会显示错误
```

### 更新go.mod中的模块路径

记得将 `github.com/dengdeng-harmenyos/server` 替换为你的实际仓库路径：

```bash
cd server
go mod edit -module github.com/你的用户名/你的仓库名
```

## 进阶用法

### 多环境配置

可以为不同环境配置不同的Secrets：

```yaml
# .github/workflows/build-staging.yml
env:
  AGCONNECT_JSON: ${{ secrets.AGCONNECT_JSON_STAGING }}
  PRIVATE_JSON: ${{ secrets.PRIVATE_JSON_STAGING }}
```

### 添加更多构建选项

```yaml
- name: Build with custom flags
  run: |
    go build -ldflags "\
      -X '...embeddedAgConnectJSON=$AGCONNECT_ESCAPED' \
      -X '...embeddedPrivateJSON=$PRIVATE_ESCAPED' \
      -X 'main.Version=${{ github.ref_name }}' \
      -X 'main.CommitHash=${{ github.sha }}' \
      -s -w" \
      -o bin/dengdeng-server \
      cmd/server/main.go
```

## 参考资料

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GitHub Secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- [Go Build Options](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies)
