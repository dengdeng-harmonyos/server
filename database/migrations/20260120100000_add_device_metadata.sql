-- Migration: 20260120100000_add_device_metadata
-- Description: 添加设备元数据字段用于更详细的设备信息追踪

-- 添加设备识别字段
ALTER TABLE devices ADD COLUMN IF NOT EXISTS device_model VARCHAR(100);
ALTER TABLE devices ADD COLUMN IF NOT EXISTS device_manufacturer VARCHAR(100);

-- 添加注释
COMMENT ON COLUMN devices.device_model IS '设备型号，如 Mate60 Pro';
COMMENT ON COLUMN devices.device_manufacturer IS '设备制造商，如 HUAWEI';

-- 创建索引以优化查询
CREATE INDEX IF NOT EXISTS idx_devices_device_model ON devices(device_model);
