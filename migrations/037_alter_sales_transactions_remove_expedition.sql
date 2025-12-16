-- UP
-- Remove expedition-related columns from sales_transactions
-- These fields are now in the shippings table

-- Drop the index first
DROP INDEX IF EXISTS idx_sales_transactions_expedition_id;

-- Remove the columns
ALTER TABLE sales_transactions DROP COLUMN IF EXISTS expedition_id;
ALTER TABLE sales_transactions DROP COLUMN IF EXISTS expedition_price;
ALTER TABLE sales_transactions DROP COLUMN IF EXISTS no_resi;
