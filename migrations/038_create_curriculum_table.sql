-- UP
-- Migration: Create curriculum table
-- Description: Creates table to store curriculum data (K13, Merdeka, Nasional)

CREATE TABLE IF NOT EXISTS curriculum (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(10) UNIQUE NOT NULL,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for faster lookups
CREATE INDEX idx_curriculum_code ON curriculum(code);
CREATE INDEX idx_curriculum_name ON curriculum(name);

-- Add comments
COMMENT ON TABLE curriculum IS 'Table to store curriculum information (K13, Merdeka, Nasional)';
COMMENT ON COLUMN curriculum.id IS 'UUID primary key (auto-generated)';
COMMENT ON COLUMN curriculum.code IS 'Unique curriculum code (e.g., K13, MER, NAS)';
COMMENT ON COLUMN curriculum.name IS 'Full name of the curriculum';
COMMENT ON COLUMN curriculum.description IS 'Optional description';

-- DOWN
-- DROP TABLE IF EXISTS curriculum;
