-- UP
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sales_transaction_id UUID NOT NULL REFERENCES sales_transactions(id) ON DELETE CASCADE,
    no_payment VARCHAR(50) UNIQUE NOT NULL,
    payment_date TIMESTAMP NOT NULL,
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    note TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for faster lookups
CREATE INDEX idx_payments_sales_transaction_id ON payments(sales_transaction_id);
CREATE INDEX idx_payments_no_payment ON payments(no_payment);
CREATE INDEX idx_payments_payment_date ON payments(payment_date);
