-- Migration: Alter publishers table - add code and description fields
-- Description: Adds code (unique) and description fields to publishers table

ALTER TABLE publishers ADD COLUMN code VARCHAR(50) UNIQUE NOT NULL;
ALTER TABLE publishers ADD COLUMN description TEXT;

-- Create index for code field
CREATE INDEX IF NOT EXISTS idx_publishers_code ON publishers(code);

-- Add comments
COMMENT ON COLUMN publishers.code IS 'Unique code for the publisher';
COMMENT ON COLUMN publishers.description IS 'Description of the publisher';
