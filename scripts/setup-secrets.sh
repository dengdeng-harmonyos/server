#!/bin/bash

# Secrets 配置助手脚本
# 自动生成所有需要添加到 GitHub Secrets 的值

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
CONFIG_DIR="$PROJECT_ROOT/config"

echo "🔐 GitHub Secrets 配置助手"
echo ""

# 检查必要的工具
if ! command -v jq &> /dev/null; then
    echo "❌ 错误：未安装 jq"
    echo ""
    echo "请安装 jq："
    echo "  macOS:    brew install jq"
    echo "  Ubuntu:   sudo apt-get install jq"
    echo "  CentOS:   sudo yum install jq"
    exit 1
fi

if ! command -v openssl &> /dev/null; then
    echo "❌ 错误：未安装 openssl"
    exit 1
fi

# 检查配置文件
AGCONNECT_FILE="$CONFIG_DIR/agconnect-services.json"
PRIVATE_FILE="$CONFIG_DIR/private.json"

if [ ! -f "$AGCONNECT_FILE" ]; then
    echo "❌ 错误：找不到 agconnect-services.json"
    echo "位置：$AGCONNECT_FILE"
    echo ""
    echo "请参考 config/SERVICE_ACCOUNT_SETUP.md 获取配置文件"
    exit 1
fi

if [ ! -f "$PRIVATE_FILE" ]; then
    echo "❌ 错误：找不到 private.json"
    echo "位置：$PRIVATE_FILE"
    echo ""
    echo "请参考 config/SERVICE_ACCOUNT_SETUP.md 获取配置文件"
    exit 1
fi

# 验证 JSON 格式
echo "🔍 验证配置文件格式..."
if ! jq -e . "$AGCONNECT_FILE" >/dev/null 2>&1; then
    echo "❌ 错误：agconnect-services.json 格式无效"
    exit 1
fi

if ! jq -e . "$PRIVATE_FILE" >/dev/null 2>&1; then
    echo "❌ 错误：private.json 格式无效"
    exit 1
fi

echo "✅ 配置文件格式正确"
echo ""

# 生成加密密钥
echo "🔑 生成 Push Token 加密密钥..."
ENCRYPTION_KEY=$(openssl rand -base64 24)
echo "✅ 加密密钥已生成"
echo ""

# 压缩 JSON（用于 GitHub Secrets）
AGCONNECT_JSON=$(cat "$AGCONNECT_FILE" | jq -c .)
PRIVATE_JSON=$(cat "$PRIVATE_FILE" | jq -c .)

# 输出配置信息
echo "═══════════════════════════════════════════════════════════════"
echo "📋 GitHub Secrets 配置值"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "将以下内容添加到 GitHub 仓库的 Secrets："
echo "路径：Settings → Secrets and variables → Actions → New repository secret"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Secret 1 of 3"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Name: AGCONNECT_JSON"
echo ""
echo "Value:"
echo "$AGCONNECT_JSON"
echo ""
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Secret 2 of 3"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Name: PRIVATE_JSON"
echo ""
echo "Value:"
echo "$PRIVATE_JSON"
echo ""
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Secret 3 of 3"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Name: PUSH_TOKEN_ENCRYPTION_KEY"
echo ""
echo "Value:"
echo "$ENCRYPTION_KEY"
echo ""
echo ""

# 创建本地 .env 文件
echo "📝 创建本地 .env 文件..."
cat > "$PROJECT_ROOT/.env" <<EOF
# Push Token 加密密钥（自动生成）
# 注意：此密钥也应添加到 GitHub Secrets 中
PUSH_TOKEN_ENCRYPTION_KEY=$ENCRYPTION_KEY

# 如果需要，你可以添加其他环境变量
# 例如：
# GIN_MODE=debug
# PORT=8080
EOF

echo "✅ 已创建 $PROJECT_ROOT/.env"
echo ""

# 输出后续步骤
echo "═══════════════════════════════════════════════════════════════"
echo "✅ 配置生成完成！"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "📚 后续步骤："
echo ""
echo "1️⃣  添加 Secrets 到 GitHub："
echo "   - 打开你的仓库"
echo "   - Settings → Secrets and variables → Actions"
echo "   - 点击 'New repository secret'"
echo "   - 将上面的3个 Name/Value 对依次添加"
echo ""
echo "2️⃣  验证 .gitignore 配置："
echo "   - 确保 .env 文件不会被提交"
echo "   - 确保 config/*.json 文件不会被提交"
echo ""
echo "3️⃣  提交代码触发构建："
echo "   $ git add ."
echo "   $ git commit -m \"feat: Configure GitHub Actions with secrets\""
echo "   $ git push origin main"
echo ""
echo "4️⃣  查看构建结果："
echo "   - 进入 GitHub 仓库的 Actions 标签"
echo "   - 查看最新的 workflow run"
echo "   - 验证构建成功"
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo ""

# 输出安全提示
echo "🔒 安全提示："
echo "   ⚠️  不要将 .env 文件提交到 Git"
echo "   ⚠️  不要将 config/*.json 文件提交到 Git"
echo "   ⚠️  定期轮换加密密钥（建议每季度一次）"
echo "   ⚠️  只授权信任的人访问 GitHub Secrets"
echo ""

# 检查 .gitignore
if [ -f "$PROJECT_ROOT/.gitignore" ]; then
    if ! grep -q "^\.env$" "$PROJECT_ROOT/.gitignore"; then
        echo "⚠️  警告：.gitignore 中未包含 .env"
        echo "建议添加：echo '.env' >> .gitignore"
        echo ""
    fi
    if ! grep -q "config/.*\.json" "$PROJECT_ROOT/.gitignore"; then
        echo "⚠️  警告：.gitignore 中未包含 config/*.json"
        echo "建议添加：echo 'config/*.json' >> .gitignore"
        echo ""
    fi
fi

echo "💡 提示：保存此输出到安全的地方，以备将来参考"
echo ""
