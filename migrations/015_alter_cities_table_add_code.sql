-- Migration: Alter cities table - add code field
-- Description: Adds unique code field to cities table

ALTER TABLE cities ADD COLUMN code VARCHAR(50) UNIQUE NOT NULL;

-- Create index for code field
CREATE INDEX IF NOT EXISTS idx_cities_code ON cities(code);

-- Add comment
COMMENT ON COLUMN cities.code IS 'Unique code for the city';
