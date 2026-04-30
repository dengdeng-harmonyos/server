-- Migration: 005_app_update_releases
-- Description: Track App update release history seeded by image manifest
-- Date: 2026-04-30

CREATE TABLE IF NOT EXISTS app_update_releases (
    platform VARCHAR(32) NOT NULL DEFAULT 'harmonyos',
    version_code BIGINT NOT NULL,
    version_name VARCHAR(64) NOT NULL DEFAULT '',
    min_version_code BIGINT NOT NULL DEFAULT 0,
    force_update BOOLEAN NOT NULL DEFAULT TRUE,
    store_url TEXT NOT NULL DEFAULT 'store://appgallery.huawei.com/app/detail?id=top.yidingyaojizhu.dengdeng',
    release_notes TEXT NOT NULL DEFAULT '',
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    source VARCHAR(64) NOT NULL DEFAULT 'manifest',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (platform, version_code)
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_app_update_releases_updated_at ON app_update_releases;
CREATE TRIGGER update_app_update_releases_updated_at
    BEFORE UPDATE ON app_update_releases
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE app_update_releases IS 'App update release history seeded from image manifest.';
COMMENT ON COLUMN app_update_releases.platform IS 'Client platform, currently harmonyos.';
COMMENT ON COLUMN app_update_releases.version_code IS 'Released App versionCode.';
COMMENT ON COLUMN app_update_releases.min_version_code IS 'Minimum versionCode allowed when this release is active.';
COMMENT ON COLUMN app_update_releases.enabled IS 'When true, startup seed may activate this release as current policy.';
