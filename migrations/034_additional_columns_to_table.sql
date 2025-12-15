-- Migration: additional_columns_to_table
-- Description: [Add description here]
-- Created: 2025-12-14 19:40:06

-- Ensure UUID extension is enabled (if needed)
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Modify the books table
ALTER TABLE books ADD COLUMN merk_buku_id UUID REFERENCES merk_buku(id) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE books ADD COLUMN periode INTEGER DEFAULT 1;

-- Modify the sales_associates table
ALTER TABLE sales_associates ADD COLUMN no_ktp VARCHAR(50);

-- Modify the sales_transactions table
ALTER TABLE sales_transactions ADD COLUMN biller_id UUID REFERENCES billers(id) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE sales_transactions ADD COLUMN no_resi VARCHAR(50);

-- Modify the sales_transaction_installments table
ALTER TABLE sales_transaction_installments ADD COLUMN no_installment VARCHAR(100) UNIQUE NOT NULL;

-- Add indexes if needed
CREATE INDEX IF NOT EXISTS idx_books_merk_buku_id ON books(merk_buku_id);
CREATE INDEX IF NOT EXISTS idx_sales_transactions_biller_id ON sales_transactions(biller_id);
CREATE INDEX IF NOT EXISTS idx_sales_transaction_installments_no_installment ON sales_transaction_installments(no_installment);

-- Add comments to table and columns
COMMENT ON COLUMN books.merk_buku_id IS 'Reference to merk_buku table (book brand)';
COMMENT ON COLUMN books.periode IS 'Semester book available (default: 1)';
COMMENT ON COLUMN sales_associates.no_ktp IS 'KTP number of the book (optional)';
COMMENT ON COLUMN sales_transactions.no_resi IS 'Expedition receipt number (optional)';
COMMENT ON COLUMN sales_transactions.biller_id IS 'Reference to billers table (optional)';
