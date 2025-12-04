-- Migration: Create jenis_buku table
-- Description: Creates table to store book type data

CREATE TABLE IF NOT EXISTS jenis_buku (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_jenis_buku_id ON jenis_buku(id);
CREATE INDEX IF NOT EXISTS idx_jenis_buku_name ON jenis_buku(name);

-- Add comments to table and columns
COMMENT ON TABLE jenis_buku IS 'Table to store book type information';
COMMENT ON COLUMN jenis_buku.id IS 'UUID primary key (auto-generated)';
COMMENT ON COLUMN jenis_buku.name IS 'Name of the book type (unique)';
COMMENT ON COLUMN jenis_buku.created_at IS 'Timestamp of record creation';
COMMENT ON COLUMN jenis_buku.updated_at IS 'Timestamp of last update';
