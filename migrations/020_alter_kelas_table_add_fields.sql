-- Migration: Alter kelas table - add code and description fields
-- Description: Adds code (unique) and description fields to kelas table

ALTER TABLE kelas ADD COLUMN code VARCHAR(50) UNIQUE NOT NULL;
ALTER TABLE kelas ADD COLUMN description TEXT;

-- Create index for code field
CREATE INDEX IF NOT EXISTS idx_kelas_code ON kelas(code);

-- Add comments
COMMENT ON COLUMN kelas.code IS 'Unique code for the kelas';
COMMENT ON COLUMN kelas.description IS 'Description of the kelas';
