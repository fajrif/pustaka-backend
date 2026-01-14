-- UP
-- Migration: Create purchase_transactions and purchase_transaction_items tables
-- Description: System for tracking book purchases from suppliers (publishers)
--   - Uses existing publishers table as suppliers
--   - Stock increases only when status = 1 (completed)
--   - Supports receipt image upload

CREATE TABLE IF NOT EXISTS purchase_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    supplier_id UUID NOT NULL REFERENCES publishers(id) ON DELETE RESTRICT,
    no_invoice VARCHAR(50) UNIQUE NOT NULL,
    purchase_date TIMESTAMP NOT NULL,
    total_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    status INTEGER NOT NULL DEFAULT 0 CHECK (status IN (0, 1, 2)),
    receipt_image_url TEXT,
    note TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Status values:
-- 0 = pending (draft, stock not affected)
-- 1 = completed (stock increased)
-- 2 = cancelled (stock not affected)

CREATE TABLE IF NOT EXISTS purchase_transaction_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    purchase_transaction_id UUID NOT NULL REFERENCES purchase_transactions(id) ON DELETE CASCADE,
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price NUMERIC(15, 2) NOT NULL,
    subtotal NUMERIC(15, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for faster lookups
CREATE INDEX idx_purchase_transactions_supplier_id ON purchase_transactions(supplier_id);
CREATE INDEX idx_purchase_transactions_no_invoice ON purchase_transactions(no_invoice);
CREATE INDEX idx_purchase_transactions_purchase_date ON purchase_transactions(purchase_date);
CREATE INDEX idx_purchase_transactions_status ON purchase_transactions(status);

CREATE INDEX idx_purchase_transaction_items_purchase_id ON purchase_transaction_items(purchase_transaction_id);
CREATE INDEX idx_purchase_transaction_items_book_id ON purchase_transaction_items(book_id);

-- Add comments
COMMENT ON TABLE purchase_transactions IS 'Table to store book purchase transactions from suppliers';
COMMENT ON COLUMN purchase_transactions.id IS 'UUID primary key (auto-generated)';
COMMENT ON COLUMN purchase_transactions.supplier_id IS 'Reference to publishers table (supplier)';
COMMENT ON COLUMN purchase_transactions.no_invoice IS 'Unique invoice number (auto-generated with PRC prefix)';
COMMENT ON COLUMN purchase_transactions.purchase_date IS 'Date of purchase';
COMMENT ON COLUMN purchase_transactions.total_amount IS 'Total amount of the purchase';
COMMENT ON COLUMN purchase_transactions.status IS 'Status: 0=pending, 1=completed, 2=cancelled';
COMMENT ON COLUMN purchase_transactions.receipt_image_url IS 'Optional receipt image URL';
COMMENT ON COLUMN purchase_transactions.note IS 'Optional notes';

COMMENT ON TABLE purchase_transaction_items IS 'Table to store items in purchase transactions';
COMMENT ON COLUMN purchase_transaction_items.purchase_transaction_id IS 'Reference to purchase_transactions table';
COMMENT ON COLUMN purchase_transaction_items.book_id IS 'Reference to books table';
COMMENT ON COLUMN purchase_transaction_items.quantity IS 'Quantity purchased';
COMMENT ON COLUMN purchase_transaction_items.price IS 'Unit price at purchase time';
COMMENT ON COLUMN purchase_transaction_items.subtotal IS 'Subtotal (price * quantity)';

-- DOWN
-- DROP TABLE IF EXISTS purchase_transaction_items;
-- DROP TABLE IF EXISTS purchase_transactions;
