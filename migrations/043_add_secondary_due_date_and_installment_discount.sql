-- UP
ALTER TABLE sales_transactions ADD COLUMN secondary_due_date TIMESTAMP;

ALTER TABLE sales_transaction_installments ADD COLUMN discount_percentage NUMERIC(5,2) NOT NULL DEFAULT 0;
ALTER TABLE sales_transaction_installments ADD COLUMN discount_amount NUMERIC(15,2) NOT NULL DEFAULT 0;
