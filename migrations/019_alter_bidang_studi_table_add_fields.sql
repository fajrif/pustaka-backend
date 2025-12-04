-- Migration: Alter bidang_studi table - add code and description fields
-- Description: Adds code (unique) and description fields to bidang_studi table

ALTER TABLE bidang_studi ADD COLUMN code VARCHAR(50) UNIQUE NOT NULL;
ALTER TABLE bidang_studi ADD COLUMN description TEXT;

-- Create index for code field
CREATE INDEX IF NOT EXISTS idx_bidang_studi_code ON bidang_studi(code);

-- Add comments
COMMENT ON COLUMN bidang_studi.code IS 'Unique code for the bidang studi';
COMMENT ON COLUMN bidang_studi.description IS 'Description of the bidang studi';
