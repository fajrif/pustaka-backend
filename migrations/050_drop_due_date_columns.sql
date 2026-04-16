-- UP
ALTER TABLE sales_transactions DROP COLUMN IF EXISTS due_date;
ALTER TABLE sales_transactions DROP COLUMN IF EXISTS secondary_due_date;
