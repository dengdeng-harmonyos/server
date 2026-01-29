#!/bin/bash

# 本地构建脚本 - 从配置文件读取并注入到编译过程
# 用于开发和测试环境

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
CONFIG_DIR="$PROJECT_ROOT/config"

# 配置文件路径
AGCONNECT_FILE="$CONFIG_DIR/agconnect-services.json"
PRIVATE_FILE="$CONFIG_DIR/private.json"

echo "🔨 Building dengdeng-server with embedded secrets..."

# 检查配置文件是否存在
if [ ! -f "$AGCONNECT_FILE" ]; then
    echo "❌ Error: agconnect-services.json not found in $CONFIG_DIR"
    exit 1
fi

if [ ! -f "$PRIVATE_FILE" ]; then
    echo "❌ Error: private.json not found in $CONFIG_DIR"
    exit 1
fi

# 读取JSON文件并转义
echo "📖 Reading configuration files..."
AGCONNECT_JSON=$(cat "$AGCONNECT_FILE" | jq -c . | sed 's/"/\\"/g')
PRIVATE_JSON=$(cat "$PRIVATE_FILE" | jq -c . | sed 's/"/\\"/g')

# 生成或读取加密密钥
if [ -z "$PUSH_TOKEN_ENCRYPTION_KEY" ]; then
    # 如果环境变量未设置，尝试从.env读取
    if [ -f "$PROJECT_ROOT/.env" ]; then
        ENCRYPTION_KEY=$(grep "^PUSH_TOKEN_ENCRYPTION_KEY=" "$PROJECT_ROOT/.env" | cut -d'=' -f2-)
    fi
    
    # 如果还是没有，生成一个随机密钥
    if [ -z "$ENCRYPTION_KEY" ]; then
        echo "⚠️  Warning: PUSH_TOKEN_ENCRYPTION_KEY not set, generating random key..."
        ENCRYPTION_KEY=$(openssl rand -base64 24)
        echo "🔑 Generated encryption key: $ENCRYPTION_KEY"
        echo "💾 Save this key to .env file: PUSH_TOKEN_ENCRYPTION_KEY=$ENCRYPTION_KEY"
    fi
else
    ENCRYPTION_KEY="$PUSH_TOKEN_ENCRYPTION_KEY"
    echo "🔑 Using encryption key from environment"
fi

# 构建二进制文件
echo "🔧 Compiling with ldflags..."
cd "$PROJECT_ROOT"

go build -ldflags "\
  -X 'github.com/dengdeng-harmonyos/server/internal/config.embeddedAgConnectJSON=$AGCONNECT_JSON' \
  -X 'github.com/dengdeng-harmonyos/server/internal/config.embeddedPrivateJSON=$PRIVATE_JSON' \
  -X 'github.com/dengdeng-harmonyos/server/internal/config.embeddedEncryptionKey=$ENCRYPTION_KEY' \
  -s -w" \
  -o bin/dengdeng-server \
  cmd/server/main.go

echo "✅ Build completed: bin/dengdeng-server"

# 检查二进制文件大小
SIZE=$(du -h bin/dengdeng-server | cut -f1)
echo "📦 Binary size: $SIZE"

echo ""
echo "🚀 You can now run: ./bin/dengdeng-server"
