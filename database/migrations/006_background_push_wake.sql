-- Add per-device throttle state for low-frequency background sync wake pushes.

ALTER TABLE devices
ADD COLUMN IF NOT EXISTS last_background_push_attempt_at TIMESTAMPTZ;

COMMENT ON COLUMN devices.last_background_push_attempt_at IS '最近一次后台唤醒Push尝试时间 (TIMESTAMPTZ, UTC)';
