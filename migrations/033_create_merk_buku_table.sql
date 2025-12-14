-- Migration: Create merk_buku table (Version 3 - with code and description)
-- Description: Creates table to store book brand data with UUID, code and description

-- Create table merk_buku
CREATE TABLE IF NOT EXISTS merk_buku (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    bantuan_promosi INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_merk_buku_id ON merk_buku(id);
CREATE INDEX IF NOT EXISTS idx_merk_buku_code ON merk_buku(code);
CREATE INDEX IF NOT EXISTS idx_merk_buku_name ON merk_buku(name);

-- Add comments to table and columns
COMMENT ON TABLE merk_buku IS 'Table to store book brand information';
COMMENT ON COLUMN merk_buku.id IS 'UUID primary key (auto-generated)';
COMMENT ON COLUMN merk_buku.code IS 'Unique code for book brand';
COMMENT ON COLUMN merk_buku.name IS 'Name of the book brand';
COMMENT ON COLUMN merk_buku.description IS 'Description of the book brand';
COMMENT ON COLUMN merk_buku.bantuan_promosi IS 'Promotional assistance flag (0 or 1)';
COMMENT ON COLUMN merk_buku.created_at IS 'Timestamp of record creation';
COMMENT ON COLUMN merk_buku.updated_at IS 'Timestamp of last update'
