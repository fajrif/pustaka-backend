-- Migration: Create merk_buku table (Version 2 - with UUID user_id)
-- Description: Creates table to store book brand/publisher data with UUID foreign key to users

-- Create table merk_buku
CREATE TABLE IF NOT EXISTS merk_buku (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kode_merk VARCHAR(10) UNIQUE NOT NULL,
    nama_merk VARCHAR(100) NOT NULL,
    bantuan_promosi INTEGER DEFAULT 0,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    tstamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_merk_buku_id ON merk_buku(id);
CREATE INDEX IF NOT EXISTS idx_merk_buku_kode_merk ON merk_buku(kode_merk);
CREATE INDEX IF NOT EXISTS idx_merk_buku_nama_merk ON merk_buku(nama_merk);
CREATE INDEX IF NOT EXISTS idx_merk_buku_user_id ON merk_buku(user_id);
CREATE INDEX IF NOT EXISTS idx_merk_buku_tstamp ON merk_buku(tstamp);

-- Add comments to table and columns
COMMENT ON TABLE merk_buku IS 'Table to store book brand/publisher information';
COMMENT ON COLUMN merk_buku.id IS 'UUID primary key (auto-generated)';
COMMENT ON COLUMN merk_buku.kode_merk IS 'Unique code for book brand';
COMMENT ON COLUMN merk_buku.nama_merk IS 'Name of the book brand';
COMMENT ON COLUMN merk_buku.bantuan_promosi IS 'Promotional assistance flag (0 or 1)';
COMMENT ON COLUMN merk_buku.user_id IS 'UUID reference to users table (who created/modified)';
COMMENT ON COLUMN merk_buku.tstamp IS 'Timestamp of record creation/modification';
