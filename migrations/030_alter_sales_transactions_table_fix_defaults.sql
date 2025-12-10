-- UP
-- Drop the old check constraint
ALTER TABLE sales_transactions DROP CONSTRAINT IF EXISTS sales_transactions_status_check;

-- Add new check constraint with status 0, 1, 2
ALTER TABLE sales_transactions ADD CONSTRAINT sales_transactions_status_check CHECK (status IN (0, 1, 2));

-- Alter the payment_type to have default value 'T'
ALTER TABLE sales_transactions ALTER COLUMN payment_type SET DEFAULT 'T';

-- Alter the status to have default value 0
ALTER TABLE sales_transactions ALTER COLUMN status SET DEFAULT 0;
