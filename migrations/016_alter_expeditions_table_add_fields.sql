-- Migration: Alter expeditions table - add code and description fields
-- Description: Adds code (unique) and description fields to expeditions table

ALTER TABLE expeditions ADD COLUMN code VARCHAR(50) UNIQUE NOT NULL;
ALTER TABLE expeditions ADD COLUMN description TEXT;

-- Create index for code field
CREATE INDEX IF NOT EXISTS idx_expeditions_code ON expeditions(code);

-- Add comments
COMMENT ON COLUMN expeditions.code IS 'Unique code for the expedition';
COMMENT ON COLUMN expeditions.description IS 'Description of the expedition';
