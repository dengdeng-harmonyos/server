-- 数据库初始化脚本
-- 用于手动创建数据库表

-- 设备表（简化版）
CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    device_key VARCHAR(64) UNIQUE NOT NULL,           -- 服务端生成的随机UUID
    push_token TEXT NOT NULL UNIQUE,                  -- 加密后的华为Push Token
    device_type VARCHAR(50),                          -- phone/tablet/watch
    os_version VARCHAR(50),                           -- HarmonyOS版本
    app_version VARCHAR(50),                          -- App版本
    is_active BOOLEAN DEFAULT TRUE,                   -- 是否活跃
    last_active_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 推送统计表（仅统计数据，不记录内容）
CREATE TABLE IF NOT EXISTS push_statistics (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,                               -- 统计日期
    push_type VARCHAR(20) NOT NULL,                   -- notification/form/background
    total_count INTEGER DEFAULT 0,                    -- 总推送数
    success_count INTEGER DEFAULT 0,                  -- 成功数
    failed_count INTEGER DEFAULT 0,                   -- 失败数
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_date_type UNIQUE(date, push_type)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_devices_device_key ON devices(device_key);
CREATE INDEX IF NOT EXISTS idx_devices_is_active ON devices(is_active);
CREATE INDEX IF NOT EXISTS idx_devices_last_active ON devices(last_active_at);
CREATE INDEX IF NOT EXISTS idx_push_stats_date ON push_statistics(date);
CREATE INDEX IF NOT EXISTS idx_push_stats_type ON push_statistics(push_type);

-- 清理不活跃设备的定时任务（可选）
-- 建议设置数据库定时任务，每月执行一次
-- DELETE FROM devices WHERE is_active = FALSE AND updated_at < NOW() - INTERVAL '90 days';

-- 插入测试数据（可选）
-- INSERT INTO devices (device_key, push_token, device_type, os_version, app_version)
-- VALUES ('test-device-key-123', 'encrypted_token_here', 'phone', 'HarmonyOS 5.0', '1.0.0');
