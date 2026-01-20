# 敏感配置管理 - 快速开始

## 📋 概述

项目已配置通过GitHub Actions Secrets安全地管理敏感配置，实现：
- ✅ 源代码不包含敏感信息，可以安全开源
- ✅ 编译时从GitHub Secrets注入配置
- ✅ 灵活配置不同环境
- ✅ 编译后的二进制文件包含配置，无需外部文件

## 🚀 快速开始

### 本地开发

1. **准备配置文件**
   ```bash
   cd server/config
   # 放置你的配置文件
   # - agconnect-services.json
   # - private.json
   ```

2. **构建**
   ```bash
   cd server
   ./scripts/build-with-secrets.sh
   ```

3. **运行**
   ```bash
   ./bin/dengdeng-server
   ```

### GitHub Actions自动构建

1. **配置Secrets**
   
   进入仓库 Settings → Secrets and variables → Actions，添加：
   
   - `AGCONNECT_JSON`: agconnect-services.json的内容（压缩为一行）
   - `PRIVATE_JSON`: private.json的内容（压缩为一行）
   
   获取压缩内容：
   ```bash
   cat config/agconnect-services.json | jq -c .
   cat config/private.json | jq -c .
   ```

2. **触发构建**
   
   推送代码或创建tag：
   ```bash
   git push origin main
   # 或
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. **下载产物**
   
   - Actions页面下载Artifacts（保留30天）
   - 或从Releases下载（tag触发时）

## 📁 文件结构

```
server/
├── .github/
│   └── workflows/
│       └── build.yml              # GitHub Actions工作流
├── config/
│   ├── README.md                  # 配置文件说明
│   ├── agconnect-services.json    # （本地，不提交）
│   └── private.json               # （本地，不提交）
├── internal/
│   └── config/
│       ├── config.go              # 配置加载逻辑
│       └── embedded_secrets.go    # 编译时注入点
├── scripts/
│   └── build-with-secrets.sh      # 本地构建脚本
├── GITHUB_ACTIONS_SETUP.md        # 详细配置指南
└── SECRETS_QUICKSTART.md          # 本文件
```

## 🔒 安全性

### 已实现
- ✅ 配置文件已加入 `.gitignore`
- ✅ GitHub Secrets加密存储
- ✅ 编译时注入，不在源码中
- ✅ 构建日志不输出secrets内容

### 注意事项
- ⚠️ 不要在代码中打印配置内容
- ⚠️ 二进制文件仍可能被反编译
- ⚠️ 定期轮换密钥
- ⚠️ 限制仓库访问权限

## 🛠️ 常用命令

### 本地编译
```bash
# 使用脚本（推荐）
./scripts/build-with-secrets.sh

# 手动编译
AGCONNECT_JSON=$(cat config/agconnect-services.json | jq -c . | sed 's/"/\\"/g')
PRIVATE_JSON=$(cat config/private.json | jq -c . | sed 's/"/\\"/g')

go build -ldflags "\
  -X 'github.com/dengdeng-harmenyos/server/internal/config.embeddedAgConnectJSON=$AGCONNECT_JSON' \
  -X 'github.com/dengdeng-harmenyos/server/internal/config.embeddedPrivateJSON=$PRIVATE_JSON' \
  -s -w" \
  -o bin/dengdeng-server \
  cmd/server/main.go
```

### 验证配置
```bash
# 检查JSON格式
jq . config/agconnect-services.json
jq . config/private.json

# 验证二进制已包含配置
strings bin/dengdeng-server | grep -i project_id
```

### 准备GitHub Secrets
```bash
# 生成压缩的JSON（复制输出到GitHub Secrets）
echo "AGCONNECT_JSON:"
cat config/agconnect-services.json | jq -c .

echo ""
echo "PRIVATE_JSON:"
cat config/private.json | jq -c .
```

## 📖 详细文档

- [GITHUB_ACTIONS_SETUP.md](GITHUB_ACTIONS_SETUP.md) - GitHub Actions详细配置指南
- [config/README.md](config/README.md) - 配置文件说明
- [EMBEDDED_SECRETS.md](EMBEDDED_SECRETS.md) - 嵌入配置技术说明

## ❓ 故障排除

### 构建失败
```bash
# 检查配置文件是否存在
ls -la config/

# 验证JSON格式
jq . config/agconnect-services.json
jq . config/private.json

# 检查构建脚本权限
chmod +x scripts/build-with-secrets.sh
```

### GitHub Actions失败
1. 检查Secrets是否正确配置
2. 验证JSON格式（使用jq）
3. 查看Actions日志中的具体错误
4. 确认go.mod中的模块路径正确

### 配置未生效
```bash
# 确认配置已嵌入
strings bin/dengdeng-server | grep project_id

# 检查环境变量优先级
# 环境变量 > 嵌入配置 > 配置文件
```

## 🔄 更新配置

### 本地开发
直接修改配置文件，然后重新编译：
```bash
./scripts/build-with-secrets.sh
```

### 生产环境
1. 更新GitHub Secrets中的值
2. 推送代码触发新构建
3. 下载新的二进制文件

## 📞 支持

遇到问题？查看：
1. [GitHub Actions构建日志](../../actions)
2. [Issues](../../issues)
3. 详细文档：GITHUB_ACTIONS_SETUP.md
