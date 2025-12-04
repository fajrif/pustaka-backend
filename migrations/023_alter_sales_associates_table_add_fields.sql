-- Migration: Alter sales_associates table - add code and description fields
-- Description: Adds code (unique) and description fields to sales_associates table

ALTER TABLE sales_associates ADD COLUMN code VARCHAR(50) UNIQUE NOT NULL;
ALTER TABLE sales_associates ADD COLUMN description TEXT;

-- Create index for code field
CREATE INDEX IF NOT EXISTS idx_sales_associates_code ON sales_associates(code);

-- Add comments
COMMENT ON COLUMN sales_associates.code IS 'Unique code for the sales associate';
COMMENT ON COLUMN sales_associates.description IS 'Description of the sales associate';
