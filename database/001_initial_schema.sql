-- 数据库初始化脚本 v1.0.0
-- 创建所有必需的表和索引

-- 设备表
CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(64) UNIQUE NOT NULL,            -- 服务端生成的随机UUID
    push_token TEXT NOT NULL UNIQUE,                  -- 加密后的华为Push Token
    public_key TEXT,                                  -- RSA公钥(PEM格式)
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

-- 待发送消息表（加密消息存储）
CREATE TABLE IF NOT EXISTS pending_messages (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(64) NOT NULL,
    server_name VARCHAR(255) NOT NULL,
    encrypted_aes_key TEXT NOT NULL,                  -- RSA加密的AES密钥
    encrypted_content TEXT NOT NULL,                  -- AES加密的消息内容
    iv TEXT NOT NULL,                                 -- AES IV向量
    notification_sent BOOLEAN DEFAULT FALSE,
    delivered BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    confirmed_at TIMESTAMP,
    
    CONSTRAINT fk_device_id FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE CASCADE
);

-- 数据库迁移版本表
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(14) PRIMARY KEY,                  -- 格式：YYYYMMDDHHMMSS
    description TEXT NOT NULL,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_devices_device_id ON devices(device_id);
CREATE INDEX IF NOT EXISTS idx_devices_is_active ON devices(is_active);
CREATE INDEX IF NOT EXISTS idx_devices_last_active ON devices(last_active_at);
CREATE INDEX IF NOT EXISTS idx_push_stats_date ON push_statistics(date);
CREATE INDEX IF NOT EXISTS idx_push_stats_type ON push_statistics(push_type);
CREATE INDEX IF NOT EXISTS idx_pending_device_id ON pending_messages(device_id);
CREATE INDEX IF NOT EXISTS idx_pending_delivered ON pending_messages(delivered);
CREATE INDEX IF NOT EXISTS idx_pending_expires ON pending_messages(expires_at);

-- 创建自动清理过期消息的函数
CREATE OR REPLACE FUNCTION clean_expired_messages() RETURNS void AS $$
BEGIN
    DELETE FROM pending_messages 
    WHERE expires_at < NOW() 
       OR (delivered = true AND confirmed_at < NOW() - INTERVAL '24 hours');
END;
$$ LANGUAGE plpgsql;

-- 创建更新时间戳的触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 为devices表添加更新时间戳触发器
DROP TRIGGER IF EXISTS update_devices_updated_at ON devices;
CREATE TRIGGER update_devices_updated_at
    BEFORE UPDATE ON devices
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 记录初始迁移版本
INSERT INTO schema_migrations (version, description) 
VALUES ('20260120000000', 'Initial database schema')
ON CONFLICT (version) DO NOTHING;

-- 输出初始化完成信息
DO $$
BEGIN
    RAISE NOTICE 'Database initialized successfully with version 20260120000000';
END $$;
