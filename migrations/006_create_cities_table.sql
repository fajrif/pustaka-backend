-- Migration: Create cities table
-- Description: Creates table to store city data

CREATE TABLE IF NOT EXISTS cities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_cities_id ON cities(id);
CREATE INDEX IF NOT EXISTS idx_cities_name ON cities(name);

-- Add comments to table and columns
COMMENT ON TABLE cities IS 'Table to store city information';
COMMENT ON COLUMN cities.id IS 'UUID primary key (auto-generated)';
COMMENT ON COLUMN cities.name IS 'Name of the city (unique)';
COMMENT ON COLUMN cities.created_at IS 'Timestamp of record creation';
COMMENT ON COLUMN cities.updated_at IS 'Timestamp of last update';
