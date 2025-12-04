-- Migration: Alter jenis_buku table - add code and description fields
-- Description: Adds code (unique) and description fields to jenis_buku table

ALTER TABLE jenis_buku ADD COLUMN code VARCHAR(50) UNIQUE NOT NULL;
ALTER TABLE jenis_buku ADD COLUMN description TEXT;

-- Create index for code field
CREATE INDEX IF NOT EXISTS idx_jenis_buku_code ON jenis_buku(code);

-- Add comments
COMMENT ON COLUMN jenis_buku.code IS 'Unique code for the jenis buku';
COMMENT ON COLUMN jenis_buku.description IS 'Description of the jenis buku';
