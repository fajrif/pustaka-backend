-- Migration: Create billers table
-- Description: Creates table to store biller data

CREATE TABLE IF NOT EXISTS billers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) UNIQUE,
    description TEXT,
    npwp VARCHAR(50) UNIQUE NOT NULL,
    address TEXT NOT NULL,
    city_id UUID REFERENCES cities(id) ON DELETE SET NULL ON UPDATE CASCADE,
    phone1 VARCHAR(50) NOT NULL,
    phone2 VARCHAR(50),
    fax VARCHAR(50),
    email VARCHAR(255),
    website VARCHAR(255),
    logo_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_billers_id ON billers(id);
CREATE INDEX IF NOT EXISTS idx_billers_code ON billers(code);
CREATE INDEX IF NOT EXISTS idx_billers_name ON billers(name);
CREATE INDEX IF NOT EXISTS idx_billers_npwp ON billers(npwp);
CREATE INDEX IF NOT EXISTS idx_billers_city_id ON billers(city_id);

-- Add comments to table and columns
COMMENT ON TABLE billers IS 'Table to store biller information';
COMMENT ON COLUMN billers.id IS 'UUID primary key (auto-generated)';
COMMENT ON COLUMN billers.code IS 'Unique code for the biller';
COMMENT ON COLUMN billers.name IS 'Name of the biller (unique)';
COMMENT ON COLUMN billers.description IS 'Description of the biller';
COMMENT ON COLUMN billers.npwp IS 'Tax identification number (NPWP) - unique and required';
COMMENT ON COLUMN billers.address IS 'Address of the biller';
COMMENT ON COLUMN billers.city_id IS 'Reference to cities table';
COMMENT ON COLUMN billers.phone1 IS 'Primary phone number';
COMMENT ON COLUMN billers.phone2 IS 'Secondary phone number';
COMMENT ON COLUMN billers.fax IS 'Fax number';
COMMENT ON COLUMN billers.email IS 'Email address';
COMMENT ON COLUMN billers.website IS 'Website URL';
COMMENT ON COLUMN billers.logo_url IS 'URL to the biller logo';
COMMENT ON COLUMN billers.created_at IS 'Timestamp of record creation';
COMMENT ON COLUMN billers.updated_at IS 'Timestamp of last update';
