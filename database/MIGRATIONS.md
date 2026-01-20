# 数据库迁移系统使用指南

## 概述

本项目使用自动化的数据库迁移系统，确保数据库架构始终保持最新状态。

## 目录结构

```
database/
├── 001_initial_schema.sql          # 初始数据库架构
├── migrate.sh                       # 迁移执行脚本
└── migrations/                      # 迁移文件目录
    └── 20260120100000_add_device_metadata.sql
```

## 迁移文件命名规范

```
{YYYYMMDDHHMMSS}_{description}.sql
```

示例：
- `20260120100000_add_device_metadata.sql`
- `20260121150000_create_notifications_table.sql`
- `20260122090000_add_user_preferences.sql`

### 命名要求
1. **版本号（14位）**：`YYYYMMDDHHMMSS` 格式，确保唯一性
2. **描述**：使用下划线分隔的英文描述
3. **扩展名**：必须是 `.sql`

## 工作流程

### 1. 容器启动时自动执行

每次Docker容器启动时，会自动：

1. 等待PostgreSQL就绪
2. 创建数据库（如果不存在）
3. 检查是否需要初始化（首次运行）
4. 扫描并执行所有未应用的迁移
5. 记录迁移版本到 `schema_migrations` 表

### 2. 迁移状态追踪

系统通过 `schema_migrations` 表追踪已应用的迁移：

```sql
CREATE TABLE schema_migrations (
    version VARCHAR(14) PRIMARY KEY,      -- 迁移版本号
    description TEXT NOT NULL,            -- 迁移描述
    applied_at TIMESTAMP                  -- 应用时间
);
```

### 3. 幂等性保证

- 所有迁移脚本应使用 `IF NOT EXISTS` 等语句确保幂等性
- 每个迁移在事务中执行，失败会自动回滚
- 已应用的迁移不会重复执行

## 创建新迁移

### 方法1：手动创建

1. 生成版本号（当前时间）：
   ```bash
   date +"%Y%m%d%H%M%S"
   # 输出：20260120123456
   ```

2. 创建迁移文件：
   ```bash
   touch database/migrations/20260120123456_your_description.sql
   ```

3. 编写迁移SQL：
   ```sql
   -- Migration: 20260120123456_your_description
   -- Description: 简短描述本次迁移的目的
   
   -- 你的SQL语句
   ALTER TABLE devices ADD COLUMN IF NOT EXISTS new_field TEXT;
   
   -- 添加索引
   CREATE INDEX IF NOT EXISTS idx_devices_new_field ON devices(new_field);
   
   -- 添加注释
   COMMENT ON COLUMN devices.new_field IS '字段说明';
   ```

### 方法2：使用生成脚本（推荐）

创建迁移生成工具：

```bash
#!/bin/bash
# scripts/create-migration.sh

if [ -z "$1" ]; then
    echo "Usage: ./scripts/create-migration.sh <description>"
    echo "Example: ./scripts/create-migration.sh add_user_email"
    exit 1
fi

DESCRIPTION=$1
VERSION=$(date +"%Y%m%d%H%M%S")
FILENAME="${VERSION}_${DESCRIPTION}.sql"
FILEPATH="database/migrations/${FILENAME}"

cat > "$FILEPATH" <<EOF
-- Migration: ${VERSION}_${DESCRIPTION}
-- Description: TODO: 描述本次迁移的目的

-- TODO: 添加你的SQL语句
-- 示例：
-- ALTER TABLE your_table ADD COLUMN IF NOT EXISTS new_column TEXT;
-- CREATE INDEX IF NOT EXISTS idx_your_index ON your_table(column_name);

EOF

echo "✅ Created migration file: $FILEPATH"
echo "📝 Please edit the file to add your SQL statements"
```

使用方法：
```bash
chmod +x scripts/create-migration.sh
./scripts/create-migration.sh add_user_preferences
```

## 迁移示例

### 示例1：添加新字段

```sql
-- Migration: 20260120100000_add_user_email
-- Description: 为设备表添加邮箱字段

ALTER TABLE devices ADD COLUMN IF NOT EXISTS email VARCHAR(255);

CREATE INDEX IF NOT EXISTS idx_devices_email ON devices(email);

COMMENT ON COLUMN devices.email IS '用户邮箱地址';
```

### 示例2：创建新表

```sql
-- Migration: 20260120110000_create_notifications_table
-- Description: 创建通知历史表

CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    device_key VARCHAR(64) NOT NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending',
    
    CONSTRAINT fk_notification_device FOREIGN KEY (device_key) 
        REFERENCES devices(device_key) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_notifications_device_key ON notifications(device_key);
CREATE INDEX IF NOT EXISTS idx_notifications_sent_at ON notifications(sent_at);
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
```

### 示例3：修改表结构

```sql
-- Migration: 20260120120000_modify_push_token_length
-- Description: 增加push_token字段长度

-- 注意：修改字段类型时要考虑数据兼容性
ALTER TABLE devices ALTER COLUMN push_token TYPE TEXT;

-- 如果需要，添加约束
ALTER TABLE devices ALTER COLUMN push_token SET NOT NULL;
```

### 示例4：数据迁移

```sql
-- Migration: 20260120130000_migrate_legacy_data
-- Description: 迁移旧数据到新格式

-- 更新现有数据
UPDATE devices 
SET device_type = 'phone' 
WHERE device_type IS NULL OR device_type = '';

-- 添加默认值
ALTER TABLE devices ALTER COLUMN device_type SET DEFAULT 'phone';
```

## 最佳实践

### 1. 向后兼容
- 添加字段时使用 `IF NOT EXISTS`
- 删除字段前确保代码不再使用
- 修改字段时考虑现有数据

### 2. 事务安全
- 每个迁移文件应该是一个完整的事务
- 避免在迁移中使用 `BEGIN/COMMIT`（脚本会自动处理）

### 3. 测试优先
```bash
# 在开发环境测试迁移
docker-compose down -v  # 清空数据
docker-compose up       # 重新启动并运行迁移
```

### 4. 文档化
- 在迁移文件中添加详细注释
- 说明迁移的原因和影响
- 记录相关的Issue或PR编号

### 5. 版本控制
- 所有迁移文件必须提交到Git
- 不要修改已应用的迁移文件
- 如需修改，创建新的回滚迁移

## 手动执行迁移

### 在Docker容器中
```bash
# 进入容器
docker exec -it push-server bash

# 执行迁移
/app/database/migrate.sh
```

### 在宿主机上
```bash
# 设置环境变量
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres123
export DB_NAME=push_server

# 执行迁移
./database/migrate.sh
```

## 查看迁移状态

### 查看已应用的迁移
```bash
docker exec -it push-server psql -U postgres -d push_server -c \
  "SELECT version, description, applied_at FROM schema_migrations ORDER BY version;"
```

### 查看当前数据库版本
```bash
docker exec -it push-server psql -U postgres -d push_server -c \
  "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1;"
```

## 故障排除

### 迁移失败

1. **查看日志**
   ```bash
   docker logs push-server
   ```

2. **检查迁移文件**
   ```bash
   # 验证SQL语法
   cat database/migrations/xxx.sql
   ```

3. **手动测试SQL**
   ```bash
   docker exec -it push-server psql -U postgres -d push_server
   # 然后粘贴SQL测试
   ```

### 回滚迁移

如果迁移出现问题，需要手动回滚：

```sql
-- 删除错误的迁移记录
DELETE FROM schema_migrations WHERE version = '20260120100000';

-- 回滚数据库更改（根据实际情况）
ALTER TABLE devices DROP COLUMN IF EXISTS problematic_field;
```

### 重置数据库

开发环境可以完全重置：

```bash
# 停止并删除所有数据
docker-compose down -v

# 重新启动（会重新初始化）
docker-compose up -d
```

## 生产环境注意事项

1. **备份优先**
   ```bash
   docker exec push-server pg_dump -U postgres push_server > backup.sql
   ```

2. **测试迁移**
   - 在开发环境完整测试
   - 在staging环境验证
   - 准备回滚方案

3. **监控执行**
   - 观察迁移日志
   - 检查应用状态
   - 验证数据完整性

4. **性能考虑**
   - 大表迁移可能耗时较长
   - 考虑在低峰期执行
   - 必要时临时停止服务

## 相关文件

- `database/001_initial_schema.sql` - 初始架构定义
- `database/migrate.sh` - 迁移执行脚本
- `database/migrations/` - 迁移文件目录
- `Dockerfile` - 容器启动配置
- `/start.sh` - 容器启动脚本（自动执行迁移）

## 技术细节

### schema_migrations 表结构
```sql
CREATE TABLE schema_migrations (
    version VARCHAR(14) PRIMARY KEY,
    description TEXT NOT NULL,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 迁移执行流程
1. 扫描 `database/migrations/` 目录
2. 读取已应用的迁移版本
3. 按版本号排序未应用的迁移
4. 在事务中执行每个迁移
5. 记录到 `schema_migrations` 表
6. 输出执行结果

### 错误处理
- 迁移失败时事务自动回滚
- 已应用的迁移不受影响
- 容器启动失败便于问题定位
