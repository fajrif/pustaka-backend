-- UP
CREATE TABLE IF NOT EXISTS sales_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sales_associate_id UUID NOT NULL REFERENCES sales_associates(id) ON DELETE RESTRICT,
    expedition_id UUID REFERENCES expeditions(id) ON DELETE SET NULL,
    no_invoice VARCHAR(100) UNIQUE NOT NULL,
    payment_type VARCHAR(1) NOT NULL CHECK (payment_type IN ('T', 'K')),
    transaction_date TIMESTAMP NOT NULL,
    due_date TIMESTAMP,
    expedition_price NUMERIC(15, 2) DEFAULT 0,
    total_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    status INTEGER NOT NULL DEFAULT 2 CHECK (status IN (1, 2)),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster lookups
CREATE INDEX idx_sales_transactions_sales_associate_id ON sales_transactions(sales_associate_id);
CREATE INDEX idx_sales_transactions_expedition_id ON sales_transactions(expedition_id);
CREATE INDEX idx_sales_transactions_no_invoice ON sales_transactions(no_invoice);
CREATE INDEX idx_sales_transactions_transaction_date ON sales_transactions(transaction_date);
CREATE INDEX idx_sales_transactions_status ON sales_transactions(status);
