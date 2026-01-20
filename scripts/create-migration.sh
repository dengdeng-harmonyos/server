#!/bin/bash

# 创建新的数据库迁移文件
# 用法: ./scripts/create-migration.sh <description>
# 示例: ./scripts/create-migration.sh add_user_preferences

set -e

if [ -z "$1" ]; then
    echo "❌ Error: Migration description is required"
    echo ""
    echo "Usage: $0 <description>"
    echo ""
    echo "Examples:"
    echo "  $0 add_user_email"
    echo "  $0 create_notifications_table"
    echo "  $0 modify_device_fields"
    echo ""
    exit 1
fi

DESCRIPTION=$1
VERSION=$(date +"%Y%m%d%H%M%S")
FILENAME="${VERSION}_${DESCRIPTION}.sql"
MIGRATIONS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/database/migrations"
FILEPATH="${MIGRATIONS_DIR}/${FILENAME}"

# 确保migrations目录存在
mkdir -p "$MIGRATIONS_DIR"

# 创建迁移文件
cat > "$FILEPATH" <<EOF
-- Migration: ${VERSION}_${DESCRIPTION}
-- Description: TODO: 描述本次迁移的目的和影响

-- ========================================
-- 在下面添加你的SQL语句
-- ========================================

-- 示例1: 添加新字段
-- ALTER TABLE devices ADD COLUMN IF NOT EXISTS new_field TEXT;
-- CREATE INDEX IF NOT EXISTS idx_devices_new_field ON devices(new_field);
-- COMMENT ON COLUMN devices.new_field IS '字段说明';

-- 示例2: 创建新表
-- CREATE TABLE IF NOT EXISTS new_table (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- 示例3: 修改字段
-- ALTER TABLE devices ALTER COLUMN existing_field TYPE VARCHAR(255);
-- ALTER TABLE devices ALTER COLUMN existing_field SET NOT NULL;

-- 示例4: 添加约束
-- ALTER TABLE devices ADD CONSTRAINT check_field_value 
--     CHECK (field_value IN ('value1', 'value2', 'value3'));

-- ========================================
-- 注意事项：
-- 1. 使用 IF NOT EXISTS 确保幂等性
-- 2. 考虑向后兼容性
-- 3. 大表操作注意性能
-- 4. 添加适当的索引
-- 5. 为字段添加注释说明
-- ========================================
EOF

echo "✅ Migration file created successfully!"
echo ""
echo "📄 File: $FILEPATH"
echo "🔢 Version: $VERSION"
echo "📝 Description: $DESCRIPTION"
echo ""
echo "📋 Next steps:"
echo "  1. Edit the file to add your SQL statements"
echo "  2. Test the migration locally"
echo "  3. Commit the file to version control"
echo ""
echo "🧪 Test migration:"
echo "  docker-compose down -v && docker-compose up"
echo ""
