-- Migration: Create kelas table
-- Description: Creates table to store class/grade data

CREATE TABLE IF NOT EXISTS kelas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_kelas_id ON kelas(id);
CREATE INDEX IF NOT EXISTS idx_kelas_name ON kelas(name);

-- Add comments to table and columns
COMMENT ON TABLE kelas IS 'Table to store class/grade information';
COMMENT ON COLUMN kelas.id IS 'UUID primary key (auto-generated)';
COMMENT ON COLUMN kelas.name IS 'Name of the class/grade (unique)';
COMMENT ON COLUMN kelas.created_at IS 'Timestamp of record creation';
COMMENT ON COLUMN kelas.updated_at IS 'Timestamp of last update';
