-- Migration: 002_enable_uuid
-- Description: Enable PostgreSQL pgcrypto extension for UUID generation
-- Date: 2026-02-11

-- Enable pgcrypto extension for gen_random_uuid() function
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Record migration
INSERT INTO schema_migrations (version, description) 
VALUES ('20260211000001', 'Enable UUID extensions');
