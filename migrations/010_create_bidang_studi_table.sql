-- Migration: Create bidang_studi table
-- Description: Creates table to store study field data

CREATE TABLE IF NOT EXISTS bidang_studi (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_bidang_studi_id ON bidang_studi(id);
CREATE INDEX IF NOT EXISTS idx_bidang_studi_name ON bidang_studi(name);

-- Add comments to table and columns
COMMENT ON TABLE bidang_studi IS 'Table to store study field information';
COMMENT ON COLUMN bidang_studi.id IS 'UUID primary key (auto-generated)';
COMMENT ON COLUMN bidang_studi.name IS 'Name of the study field (unique)';
COMMENT ON COLUMN bidang_studi.created_at IS 'Timestamp of record creation';
COMMENT ON COLUMN bidang_studi.updated_at IS 'Timestamp of last update';
