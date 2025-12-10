-- UP
CREATE TABLE IF NOT EXISTS sales_transaction_installments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID NOT NULL REFERENCES sales_transactions(id) ON DELETE CASCADE,
    installment_date TIMESTAMP NOT NULL,
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    note TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster lookups
CREATE INDEX idx_sales_transaction_installments_transaction_id ON sales_transaction_installments(transaction_id);
CREATE INDEX idx_sales_transaction_installments_installment_date ON sales_transaction_installments(installment_date);
