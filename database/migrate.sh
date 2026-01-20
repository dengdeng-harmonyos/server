#!/bin/bash

# 数据库迁移脚本
# 自动检测并执行所有未应用的迁移

set -e

# 配置
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres123}"
DB_NAME="${DB_NAME:-push_server}"

MIGRATIONS_DIR="/app/database/migrations"
INIT_SCHEMA="/app/database/001_initial_schema.sql"

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 等待数据库就绪
wait_for_db() {
    log_info "Waiting for PostgreSQL to be ready..."
    for i in {1..30}; do
        if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c '\q' 2>/dev/null; then
            log_info "PostgreSQL is ready!"
            return 0
        fi
        echo -n "."
        sleep 1
    done
    log_error "PostgreSQL failed to start within 30 seconds"
    return 1
}

# 创建数据库（如果不存在）
create_database() {
    log_info "Creating database if not exists..."
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "CREATE DATABASE $DB_NAME;" 2>/dev/null || true
}

# 检查是否需要初始化
check_initialization() {
    local table_exists=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -tAc "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'schema_migrations');")
    
    if [ "$table_exists" = "f" ]; then
        log_info "Database not initialized. Running initial schema..."
        if [ -f "$INIT_SCHEMA" ]; then
            PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$INIT_SCHEMA"
            log_info "Initial schema applied successfully"
        else
            log_error "Initial schema file not found: $INIT_SCHEMA"
            exit 1
        fi
    else
        log_info "Database already initialized"
    fi
}

# 获取已应用的迁移
get_applied_migrations() {
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -tAc "SELECT version FROM schema_migrations ORDER BY version;"
}

# 应用单个迁移
apply_migration() {
    local migration_file=$1
    local version=$(basename "$migration_file" .sql | cut -d'_' -f1)
    local description=$(basename "$migration_file" .sql | cut -d'_' -f2-)
    
    log_info "Applying migration: $version - $description"
    
    # 开始事务
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME <<EOF
BEGIN;

-- 执行迁移文件
\i $migration_file

-- 记录迁移版本
INSERT INTO schema_migrations (version, description) 
VALUES ('$version', '$description')
ON CONFLICT (version) DO NOTHING;

COMMIT;
EOF
    
    if [ $? -eq 0 ]; then
        log_info "Migration $version applied successfully"
        return 0
    else
        log_error "Failed to apply migration $version"
        return 1
    fi
}

# 执行迁移
run_migrations() {
    log_info "Checking for pending migrations..."
    
    if [ ! -d "$MIGRATIONS_DIR" ]; then
        log_warn "Migrations directory not found: $MIGRATIONS_DIR"
        return 0
    fi
    
    # 获取已应用的迁移
    applied_migrations=$(get_applied_migrations)
    
    # 查找所有迁移文件（按版本号排序）
    migration_files=$(find "$MIGRATIONS_DIR" -name "*.sql" -type f | sort)
    
    if [ -z "$migration_files" ]; then
        log_info "No migration files found"
        return 0
    fi
    
    local pending_count=0
    local applied_count=0
    
    for migration_file in $migration_files; do
        version=$(basename "$migration_file" .sql | cut -d'_' -f1)
        
        # 检查是否已应用
        if echo "$applied_migrations" | grep -q "^$version$"; then
            log_info "Migration $version already applied, skipping..."
            ((applied_count++))
            continue
        fi
        
        # 应用迁移
        if apply_migration "$migration_file"; then
            ((pending_count++))
        else
            log_error "Migration failed, stopping..."
            exit 1
        fi
    done
    
    log_info "Migration summary: $applied_count already applied, $pending_count newly applied"
}

# 显示当前数据库版本
show_version() {
    log_info "Current database schema version:"
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT version, description, applied_at FROM schema_migrations ORDER BY version;"
}

# 主流程
main() {
    log_info "Starting database migration process..."
    log_info "Database: $DB_NAME at $DB_HOST:$DB_PORT"
    
    # 等待数据库就绪
    wait_for_db || exit 1
    
    # 创建数据库
    create_database
    
    # 检查并执行初始化
    check_initialization
    
    # 执行迁移
    run_migrations
    
    # 显示当前版本
    show_version
    
    log_info "Database migration completed successfully!"
}

# 执行主流程
main
