-- Migration: drop_index_name_sales_associates
-- Description: [Add description here]
-- Created: 2025-12-09 19:07:32

-- Ensure UUID extension is enabled (if needed)
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Write your SQL here
ALTER TABLE sales_associates
DROP CONSTRAINT IF EXISTS sales_associates_name_key;

DROP INDEX IF EXISTS idx_sales_associates_name;

-- Add indexes if needed
-- CREATE INDEX IF NOT EXISTS idx_table_column ON table_name(column_name);

-- Add comments
-- COMMENT ON TABLE table_name IS 'Description';
-- COMMENT ON COLUMN table_name.column_name IS 'Description';
