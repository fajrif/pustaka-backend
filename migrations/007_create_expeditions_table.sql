-- Migration: Create expeditions table
-- Description: Creates table to store expedition/courier data

CREATE TABLE IF NOT EXISTS expeditions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    address TEXT NOT NULL,
    city_id UUID REFERENCES cities(id) ON DELETE SET NULL ON UPDATE CASCADE,
    area VARCHAR(255),
    phone1 VARCHAR(50) NOT NULL,
    phone2 VARCHAR(50),
    email VARCHAR(255),
    website VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_expeditions_id ON expeditions(id);
CREATE INDEX IF NOT EXISTS idx_expeditions_name ON expeditions(name);
CREATE INDEX IF NOT EXISTS idx_expeditions_city_id ON expeditions(city_id);

-- Add comments to table and columns
COMMENT ON TABLE expeditions IS 'Table to store expedition/courier information';
COMMENT ON COLUMN expeditions.id IS 'UUID primary key (auto-generated)';
COMMENT ON COLUMN expeditions.name IS 'Name of the expedition (unique)';
COMMENT ON COLUMN expeditions.address IS 'Address of the expedition';
COMMENT ON COLUMN expeditions.city_id IS 'Reference to cities table';
COMMENT ON COLUMN expeditions.area IS 'Area/region information';
COMMENT ON COLUMN expeditions.phone1 IS 'Primary phone number';
COMMENT ON COLUMN expeditions.phone2 IS 'Secondary phone number';
COMMENT ON COLUMN expeditions.email IS 'Email address';
COMMENT ON COLUMN expeditions.website IS 'Website URL';
COMMENT ON COLUMN expeditions.created_at IS 'Timestamp of record creation';
COMMENT ON COLUMN expeditions.updated_at IS 'Timestamp of last update';
