-- Migration: Create jenjang_studi table
-- Description: Creates table to store education level data

CREATE TABLE IF NOT EXISTS jenjang_studi (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_jenjang_studi_id ON jenjang_studi(id);
CREATE INDEX IF NOT EXISTS idx_jenjang_studi_name ON jenjang_studi(name);

-- Add comments to table and columns
COMMENT ON TABLE jenjang_studi IS 'Table to store education level information';
COMMENT ON COLUMN jenjang_studi.id IS 'UUID primary key (auto-generated)';
COMMENT ON COLUMN jenjang_studi.name IS 'Name of the education level (unique)';
COMMENT ON COLUMN jenjang_studi.created_at IS 'Timestamp of record creation';
COMMENT ON COLUMN jenjang_studi.updated_at IS 'Timestamp of last update';
