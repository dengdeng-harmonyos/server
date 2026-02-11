-- Migration: 003_migrate_to_uuid_keys
-- Description: Migrate all primary keys from SERIAL to UUID
-- Date: 2026-02-11
-- IMPORTANT: This migration requires pgcrypto extension (run 002_enable_uuid.sql first)

BEGIN;

-- ========================================
-- Step 1: Migrate devices table
-- ========================================

-- Drop the redundant integer id column (replaced by device_id as primary key)
ALTER TABLE devices DROP CONSTRAINT devices_pkey;
ALTER TABLE devices DROP COLUMN id;

-- Convert device_id from VARCHAR(64) to UUID type
ALTER TABLE devices 
    ALTER COLUMN device_id TYPE UUID USING device_id::uuid,
    ALTER COLUMN device_id SET DEFAULT gen_random_uuid();

-- Set device_id as the new primary key
ALTER TABLE devices ADD PRIMARY KEY (device_id);

-- Drop and recreate index (now using UUID type)
DROP INDEX IF EXISTS idx_devices_device_id;
CREATE INDEX idx_devices_device_id ON devices(device_id);

COMMENT ON COLUMN devices.device_id IS 'UUID primary key, server-generated or client-provided';


-- ========================================
-- Step 2: Migrate pending_messages table
-- ========================================

-- Drop foreign key constraint temporarily
ALTER TABLE pending_messages DROP CONSTRAINT fk_device_id;

-- Add new UUID id column
ALTER TABLE pending_messages ADD COLUMN uuid_id UUID DEFAULT gen_random_uuid();

-- Populate uuid_id for existing rows
UPDATE pending_messages SET uuid_id = gen_random_uuid() WHERE uuid_id IS NULL;

-- Convert device_id from VARCHAR to UUID
ALTER TABLE pending_messages 
    ALTER COLUMN device_id TYPE UUID USING device_id::uuid;

-- Drop old SERIAL primary key
ALTER TABLE pending_messages DROP CONSTRAINT pending_messages_pkey;
ALTER TABLE pending_messages DROP COLUMN id;

-- Rename uuid_id to id and set as primary key
ALTER TABLE pending_messages RENAME COLUMN uuid_id TO id;
ALTER TABLE pending_messages ALTER COLUMN id SET NOT NULL;
ALTER TABLE pending_messages ADD PRIMARY KEY (id);

-- Recreate foreign key constraint (now referencing UUID)
ALTER TABLE pending_messages 
    ADD CONSTRAINT fk_device_id 
    FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE CASCADE;

-- Recreate index with UUID type
DROP INDEX IF EXISTS idx_pending_device_id;
CREATE INDEX idx_pending_device_id ON pending_messages(device_id);

COMMENT ON COLUMN pending_messages.id IS 'UUID primary key';
COMMENT ON COLUMN pending_messages.device_id IS 'UUID foreign key to devices.device_id';


-- ========================================
-- Step 3: Migrate push_statistics table
-- ========================================

-- Add new UUID id column
ALTER TABLE push_statistics ADD COLUMN uuid_id UUID DEFAULT gen_random_uuid();

-- Populate uuid_id for existing rows
UPDATE push_statistics SET uuid_id = gen_random_uuid() WHERE uuid_id IS NULL;

-- Drop old SERIAL primary key
ALTER TABLE push_statistics DROP CONSTRAINT push_statistics_pkey;
ALTER TABLE push_statistics DROP COLUMN id;

-- Rename uuid_id to id and set as primary key
ALTER TABLE push_statistics RENAME COLUMN uuid_id TO id;
ALTER TABLE push_statistics ALTER COLUMN id SET NOT NULL;
ALTER TABLE push_statistics ADD PRIMARY KEY (id);

COMMENT ON COLUMN push_statistics.id IS 'UUID primary key';


-- ========================================
-- Step 4: Record migration
-- ========================================

INSERT INTO schema_migrations (version, description) 
VALUES ('20260211000002', 'Migrate primary keys to UUID type');

COMMIT;

-- Migration completed successfully
-- All tables now use UUID primary keys:
--   - devices: device_id (UUID, PRIMARY KEY)
--   - pending_messages: id (UUID), device_id (UUID, FOREIGN KEY)
--   - push_statistics: id (UUID)
