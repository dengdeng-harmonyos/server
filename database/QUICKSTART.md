# 数据库迁移快速开始

## ✅ 已完成配置

数据库迁移系统已经配置完成，每次Docker容器启动时会自动检测并应用数据库变更。

## 📁 文件结构

```
server/
├── database/
│   ├── 001_initial_schema.sql              # 初始数据库架构
│   ├── migrate.sh                          # 自动迁移脚本
│   ├── MIGRATIONS.md                       # 详细使用文档
│   └── migrations/                         # 迁移文件目录
│       └── 20260120100000_add_device_metadata.sql
├── scripts/
│   └── create-migration.sh                 # 创建迁移文件工具
└── Dockerfile                              # 已更新，支持自动迁移
```

## 🚀 快速使用

### 1. 启动容器（自动执行迁移）

```bash
cd server
docker-compose up -d
```

容器启动时会自动：
1. 初始化PostgreSQL
2. 创建数据库
3. 执行初始架构（首次）
4. 应用所有未执行的迁移

### 2. 查看迁移日志

```bash
docker logs push-server
```

你会看到类似的输出：
```
==========================================
Starting Push Server Container
==========================================
[START] Starting PostgreSQL...
[READY] PostgreSQL is ready
[MIGRATE] Running database migrations...
[INFO] Database already initialized
[INFO] Checking for pending migrations...
[INFO] Migration 20260120100000 already applied, skipping...
[MIGRATE] Database migration completed successfully
[START] Starting push server application...
```

### 3. 创建新迁移

```bash
# 使用脚本创建（推荐）
./scripts/create-migration.sh add_notification_settings

# 编辑生成的文件
vim database/migrations/20260120xxxxxx_add_notification_settings.sql
```

### 4. 测试迁移

```bash
# 完全重建（开发环境）
docker-compose down -v
docker-compose up -d

# 查看结果
docker logs push-server | grep MIGRATE
```

### 5. 查看数据库状态

```bash
# 查看已应用的迁移
docker exec -it push-server psql -U postgres -d push_server -c \
  "SELECT version, description, applied_at FROM schema_migrations ORDER BY version;"

# 查看当前版本
docker exec -it push-server psql -U postgres -d push_server -c \
  "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1;"
```

## 📝 迁移示例

### 添加新字段

```bash
# 1. 创建迁移文件
./scripts/create-migration.sh add_user_timezone

# 2. 编辑文件添加SQL
cat > database/migrations/20260120150000_add_user_timezone.sql <<'EOF'
-- Migration: 20260120150000_add_user_timezone
-- Description: 添加用户时区设置

ALTER TABLE devices ADD COLUMN IF NOT EXISTS timezone VARCHAR(50) DEFAULT 'UTC';

CREATE INDEX IF NOT EXISTS idx_devices_timezone ON devices(timezone);

COMMENT ON COLUMN devices.timezone IS '用户时区，如 Asia/Shanghai';
EOF

# 3. 重启容器应用迁移
docker-compose restart
```

## 🔍 常用命令

```bash
# 查看迁移状态
docker exec push-server /app/database/migrate.sh

# 进入数据库
docker exec -it push-server psql -U postgres -d push_server

# 查看表结构
docker exec -it push-server psql -U postgres -d push_server -c "\d devices"

# 查看所有表
docker exec -it push-server psql -U postgres -d push_server -c "\dt"

# 备份数据库
docker exec push-server pg_dump -U postgres push_server > backup_$(date +%Y%m%d).sql

# 恢复数据库
cat backup_20260120.sql | docker exec -i push-server psql -U postgres -d push_server
```

## ⚠️ 注意事项

### 开发环境
- 可以随时使用 `docker-compose down -v` 完全重置
- 测试迁移确保幂等性

### 生产环境
- **必须先备份数据库**
- 在staging环境充分测试
- 准备回滚方案
- 考虑在低峰期执行

## 📚 更多信息

详细文档请查看：
- [database/MIGRATIONS.md](database/MIGRATIONS.md) - 完整迁移文档
- [scripts/create-migration.sh](scripts/create-migration.sh) - 迁移创建工具

## 🎯 核心特性

✅ **自动执行** - 容器启动时自动检测并应用迁移  
✅ **幂等性** - 已应用的迁移不会重复执行  
✅ **事务安全** - 迁移失败自动回滚  
✅ **版本追踪** - 记录所有已应用的迁移  
✅ **向后兼容** - 使用IF NOT EXISTS等语句  

## 💡 最佳实践

1. **命名规范**
   ```
   {YYYYMMDDHHMMSS}_{description}.sql
   ```

2. **编写迁移**
   - 使用 `IF NOT EXISTS` 确保幂等性
   - 添加适当的索引和注释
   - 考虑向后兼容性

3. **测试流程**
   ```bash
   # 1. 创建迁移
   ./scripts/create-migration.sh your_migration
   
   # 2. 编辑SQL
   vim database/migrations/xxx.sql
   
   # 3. 测试
   docker-compose down -v && docker-compose up -d
   
   # 4. 验证
   docker logs push-server | grep MIGRATE
   ```

4. **版本控制**
   - 提交所有迁移文件到Git
   - 不要修改已应用的迁移
   - 需要修改时创建新的迁移

---

**迁移系统已就绪！** 🎉

现在每次重启Docker容器时，都会自动检查并应用数据库变更。
