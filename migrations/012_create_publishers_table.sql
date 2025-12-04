-- Migration: Create publishers table
-- Description: Creates table to store publisher data

CREATE TABLE IF NOT EXISTS publishers (
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
CREATE INDEX IF NOT EXISTS idx_publishers_id ON publishers(id);
CREATE INDEX IF NOT EXISTS idx_publishers_name ON publishers(name);
CREATE INDEX IF NOT EXISTS idx_publishers_city_id ON publishers(city_id);

-- Add comments to table and columns
COMMENT ON TABLE publishers IS 'Table to store publisher information';
COMMENT ON COLUMN publishers.id IS 'UUID primary key (auto-generated)';
COMMENT ON COLUMN publishers.name IS 'Name of the publisher (unique)';
COMMENT ON COLUMN publishers.address IS 'Address of the publisher';
COMMENT ON COLUMN publishers.city_id IS 'Reference to cities table';
COMMENT ON COLUMN publishers.area IS 'Area/region information';
COMMENT ON COLUMN publishers.phone1 IS 'Primary phone number';
COMMENT ON COLUMN publishers.phone2 IS 'Secondary phone number';
COMMENT ON COLUMN publishers.email IS 'Email address';
COMMENT ON COLUMN publishers.website IS 'Website URL';
COMMENT ON COLUMN publishers.created_at IS 'Timestamp of record creation';
COMMENT ON COLUMN publishers.updated_at IS 'Timestamp of last update';
