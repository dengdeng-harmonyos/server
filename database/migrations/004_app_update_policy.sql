-- Migration: 004_app_update_policy
-- Description: Store App force-update policy in database
-- Date: 2026-04-30

CREATE TABLE IF NOT EXISTS app_update_policies (
    platform VARCHAR(32) PRIMARY KEY DEFAULT 'harmonyos',
    latest_version_code BIGINT NOT NULL DEFAULT 0,
    latest_version_name VARCHAR(64) NOT NULL DEFAULT '',
    min_version_code BIGINT NOT NULL DEFAULT 0,
    force_update BOOLEAN NOT NULL DEFAULT TRUE,
    store_url TEXT NOT NULL DEFAULT 'store://appgallery.huawei.com/app/detail?id=top.yidingyaojizhu.dengdeng',
    release_notes TEXT NOT NULL DEFAULT '',
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO app_update_policies (platform)
VALUES ('harmonyos')
ON CONFLICT (platform) DO NOTHING;

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_app_update_policies_updated_at ON app_update_policies;
CREATE TRIGGER update_app_update_policies_updated_at
    BEFORE UPDATE ON app_update_policies
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE app_update_policies IS 'App force-update policy. One active policy per platform.';
COMMENT ON COLUMN app_update_policies.platform IS 'Client platform, currently harmonyos.';
COMMENT ON COLUMN app_update_policies.latest_version_code IS 'Latest published App versionCode.';
COMMENT ON COLUMN app_update_policies.min_version_code IS 'Minimum versionCode allowed to use App features.';
COMMENT ON COLUMN app_update_policies.force_update IS 'When true, any newer version forces update.';
