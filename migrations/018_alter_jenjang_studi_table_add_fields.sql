-- Migration: Alter jenjang_studi table - add code, description, and period fields
-- Description: Adds code (unique), description, and period fields to jenjang_studi table

ALTER TABLE jenjang_studi ADD COLUMN code VARCHAR(50) UNIQUE NOT NULL;
ALTER TABLE jenjang_studi ADD COLUMN description TEXT;
ALTER TABLE jenjang_studi ADD COLUMN period VARCHAR(10) DEFAULT 'S';

-- Create index for code field
CREATE INDEX IF NOT EXISTS idx_jenjang_studi_code ON jenjang_studi(code);

-- Add comments
COMMENT ON COLUMN jenjang_studi.code IS 'Unique code for the jenjang studi';
COMMENT ON COLUMN jenjang_studi.description IS 'Description of the jenjang studi';
COMMENT ON COLUMN jenjang_studi.period IS 'Period type (S/T), default: S';
