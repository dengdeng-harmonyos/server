# GitHub Actions 配置完整指南

## 📋 概述

本指南将帮助你完整配置 GitHub Actions，实现自动构建并安全地注入敏感配置。

## 🔐 需要配置的 Secrets

进入仓库 Settings → Secrets and variables → Actions，添加以下 **3个** secrets：

### 1. AGCONNECT_JSON

华为推送服务的应用配置（压缩为一行的JSON）

**获取方式：**
```bash
cd server
cat config/agconnect-services.json | jq -c .
```

复制输出的整行JSON，粘贴到 GitHub Secrets。

### 2. PRIVATE_JSON

华为OAuth 2.0服务账号私钥（压缩为一行的JSON）

**获取方式：**
```bash
cd server
cat config/private.json | jq -c .
```

复制输出的整行JSON，粘贴到 GitHub Secrets。

### 3. PUSH_TOKEN_ENCRYPTION_KEY

Push Token加密密钥（32字节字符串或Base64）

**生成方式：**
```bash
# 方法1：生成Base64密钥（推荐）
openssl rand -base64 24

# 方法2：生成32字符密钥
openssl rand -hex 16

# 方法3：使用Python生成
python3 -c "import secrets; print(secrets.token_urlsafe(32)[:32])"
```

复制生成的密钥，粘贴到 GitHub Secrets。

## 📸 配置截图参考

### 步骤1：进入 Secrets 配置页面

```
你的仓库 → Settings → Secrets and variables → Actions → New repository secret
```

### 步骤2：添加 AGCONNECT_JSON

```
Name: AGCONNECT_JSON
Secret: {"agcgw_all":{"SG":"connect-dra...（完整的压缩JSON）
```

### 步骤3：添加 PRIVATE_JSON

```
Name: PRIVATE_JSON
Secret: {"project_id":"101653523863472352",...（完整的压缩JSON）
```

### 步骤4：添加 PUSH_TOKEN_ENCRYPTION_KEY

```
Name: PUSH_TOKEN_ENCRYPTION_KEY
Secret: YourGeneratedEncryptionKey==
```

### 完成后的 Secrets 列表

你应该看到3个secrets：
- ✅ AGCONNECT_JSON
- ✅ PRIVATE_JSON
- ✅ PUSH_TOKEN_ENCRYPTION_KEY

## 📝 完整配置步骤

### 第一步：生成配置内容

在本地项目目录执行：

```bash
cd /path/to/your/project/server

# 1. 生成 AGCONNECT_JSON
echo "=== AGCONNECT_JSON ==="
cat config/agconnect-services.json | jq -c .
echo ""

# 2. 生成 PRIVATE_JSON
echo "=== PRIVATE_JSON ==="
cat config/private.json | jq -c .
echo ""

# 3. 生成 PUSH_TOKEN_ENCRYPTION_KEY
echo "=== PUSH_TOKEN_ENCRYPTION_KEY ==="
openssl rand -base64 24
echo ""
```

### 第二步：添加到 GitHub

1. 打开你的 GitHub 仓库
2. 点击 **Settings** 标签
3. 在左侧菜单找到 **Secrets and variables** → **Actions**
4. 点击 **New repository secret** 按钮
5. 依次添加上面的3个secrets

### 第三步：验证配置

提交代码触发构建：

```bash
git add .
git commit -m "feat: Configure GitHub Actions with secrets"
git push origin main
```

### 第四步：查看构建结果

1. 进入仓库的 **Actions** 标签
2. 找到最新的 workflow run
3. 查看构建日志，确认成功

预期日志输出：
```
Run go build -ldflags ...
  -X 'github.com/.../embeddedAgConnectJSON=...'
  -X 'github.com/.../embeddedPrivateJSON=...'
  -X 'github.com/.../embeddedEncryptionKey=...'
  -s -w
✓ Build completed successfully
```

## 🔧 GitHub Actions Workflow 文件

文件已创建：`.github/workflows/build.yml`

### 关键配置说明

```yaml
- name: Build with embedded secrets
  working-directory: ./server
  env:
    # 从 GitHub Secrets 读取
    AGCONNECT_JSON: ${{ secrets.AGCONNECT_JSON }}
    PRIVATE_JSON: ${{ secrets.PRIVATE_JSON }}
    ENCRYPTION_KEY: ${{ secrets.PUSH_TOKEN_ENCRYPTION_KEY }}
  run: |
    # 转义处理
    AGCONNECT_ESCAPED=$(echo "$AGCONNECT_JSON" | sed 's/"/\\"/g' | tr -d '\n')
    PRIVATE_ESCAPED=$(echo "$PRIVATE_JSON" | sed 's/"/\\"/g' | tr -d '\n')
    
    # 编译时注入
    go build -ldflags "\
      -X '...embeddedAgConnectJSON=$AGCONNECT_ESCAPED' \
      -X '...embeddedPrivateJSON=$PRIVATE_ESCAPED' \
      -X '...embeddedEncryptionKey=$ENCRYPTION_KEY' \
      -s -w" \
      -o bin/dengdeng-server \
      cmd/server/main.go
```

### 工作流触发条件

```yaml
on:
  push:
    branches: [main, master]  # 推送到主分支
    tags: ['v*']              # 创建版本标签
  pull_request:
    branches: [main, master]  # Pull Request
  workflow_dispatch:          # 手动触发
```

## 🚀 使用方法

### 自动构建（推荐）

推送代码后自动触发：

```bash
git add .
git commit -m "Your commit message"
git push origin main
```

### 手动触发

1. 进入 Actions 标签
2. 选择 "Build and Release" workflow
3. 点击 "Run workflow"
4. 选择分支，点击运行

### 版本发布

创建版本标签自动发布：

```bash
git tag v1.0.0
git push origin v1.0.0
```

会自动创建 GitHub Release 并附带二进制文件。

## 📦 产物下载

### Artifacts（保留30天）

1. 进入 Actions → 选择构建记录
2. 滚动到 "Artifacts" 部分
3. 下载 `dengdeng-server-{SHA}`

### Releases（永久）

1. 进入仓库的 Releases 页面
2. 找到对应版本
3. 下载 `dengdeng-server` 二进制文件

## 🔍 验证构建结果

下载二进制文件后验证：

```bash
# 添加执行权限
chmod +x dengdeng-server

# 测试运行（需要数据库）
./dengdeng-server

# 查看是否包含配置（不会显示完整内容，正常）
strings dengdeng-server | grep -i "project_id" | head -3
```

## 🛡️ 安全性说明

### 已实现的安全措施

✅ **Secrets 加密存储**
- GitHub 使用加密存储所有 secrets
- 只有仓库管理员可以访问

✅ **日志保护**
- Secrets 不会出现在构建日志中
- GitHub 自动屏蔽敏感信息

✅ **编译时注入**
- 配置在编译时注入到二进制
- 源代码不包含敏感信息

✅ **二进制混淆**
- 使用 `-s -w` 去除符号表
- 配置被编码在二进制中

### 最佳实践

1. **限制仓库访问**
   - 只给信任的人员仓库访问权限
   - 定期审查协作者列表

2. **定期轮换密钥**
   - 建议每季度更换 ENCRYPTION_KEY
   - 更新 Secrets 中的配置

3. **监控构建日志**
   - 检查是否有异常构建
   - 确保没有敏感信息泄露

4. **分支保护**
   ```
   Settings → Branches → Add rule
   - Require pull request reviews
   - Require status checks to pass
   ```

## ❓ 常见问题

### Q1: Secrets 配置错误怎么办？

**症状：** 构建失败，提示无法解析配置

**解决：**
1. 检查 Secrets 是否正确设置（3个都要有）
2. 验证 JSON 格式是否正确（使用 `jq` 验证）
3. 重新生成并更新 Secrets

```bash
# 验证 JSON 格式
echo "$YOUR_JSON" | jq .
```

### Q2: 如何更新 Secrets？

1. 进入 Settings → Secrets and variables → Actions
2. 点击要更新的 secret
3. 点击 "Update secret"
4. 粘贴新值，保存

### Q3: 构建成功但二进制无法运行？

**检查步骤：**

```bash
# 1. 验证文件完整性
ls -lh dengdeng-server

# 2. 检查执行权限
chmod +x dengdeng-server

# 3. 测试运行
./dengdeng-server --help

# 4. 查看详细错误
./dengdeng-server 2>&1 | head -20
```

### Q4: 如何在本地使用相同配置？

本地使用 `.env` 文件：

```bash
# 创建 .env 文件
cat > server/.env <<EOF
PUSH_TOKEN_ENCRYPTION_KEY=YourGeneratedKeyHere
EOF

# 或使用本地构建脚本
cd server
./scripts/build-with-secrets.sh
```

### Q5: 忘记了加密密钥怎么办？

**后果：**
- 无法解密已存储的 Push Token
- 所有设备需要重新注册

**解决方案：**
1. 生成新的加密密钥
2. 更新 GitHub Secrets
3. 清空数据库 devices 表
4. 通知用户重新注册设备

```sql
-- 清空设备表（慎重！）
TRUNCATE TABLE devices CASCADE;
```

## 📊 工作流程图

```
代码推送/标签创建
    ↓
GitHub Actions 触发
    ↓
读取 Secrets
  - AGCONNECT_JSON
  - PRIVATE_JSON  
  - PUSH_TOKEN_ENCRYPTION_KEY
    ↓
转义处理
    ↓
编译时注入（-ldflags -X）
    ↓
生成二进制文件
    ↓
上传 Artifact
    ↓
（如果是 tag）创建 Release
```

## 🎯 完整配置检查清单

### GitHub Secrets 配置

- [ ] 已添加 `AGCONNECT_JSON`
- [ ] 已添加 `PRIVATE_JSON`
- [ ] 已添加 `PUSH_TOKEN_ENCRYPTION_KEY`
- [ ] 所有 Secrets 格式正确
- [ ] 已测试过 JSON 格式有效

### Workflow 文件

- [ ] `.github/workflows/build.yml` 已创建
- [ ] workflow 文件语法正确
- [ ] 触发条件符合需求
- [ ] 构建步骤配置完整

### 本地环境

- [ ] `config/agconnect-services.json` 存在
- [ ] `config/private.json` 存在
- [ ] `.gitignore` 已配置忽略敏感文件
- [ ] 本地构建脚本可执行

### 安全检查

- [ ] 敏感文件未提交到 Git
- [ ] `.env` 文件在 `.gitignore` 中
- [ ] Secrets 只有管理员可访问
- [ ] 定期审查协作者权限

## 📚 相关文档

- [SECRETS_QUICKSTART.md](SECRETS_QUICKSTART.md) - Secrets 快速开始
- [GITHUB_ACTIONS_SETUP.md](GITHUB_ACTIONS_SETUP.md) - GitHub Actions 详细配置
- [DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md) - Docker 部署说明

## 🆘 获取帮助

遇到问题？

1. 查看构建日志中的详细错误信息
2. 参考本文档的常见问题部分
3. 在仓库 Issues 中搜索类似问题
4. 创建新 Issue 并附带详细信息：
   - 错误日志
   - 配置步骤
   - 环境信息

---

**配置完成后，你的构建流程将完全自动化，安全地注入所有敏感配置！** 🎉
