-- UP
CREATE TABLE IF NOT EXISTS shippings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sales_transaction_id UUID NOT NULL REFERENCES sales_transactions(id) ON DELETE CASCADE,
    expedition_id UUID NOT NULL REFERENCES expeditions(id) ON DELETE RESTRICT,
    no_resi VARCHAR(100),
    total_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for faster lookups
CREATE INDEX idx_shippings_sales_transaction_id ON shippings(sales_transaction_id);
CREATE INDEX idx_shippings_expedition_id ON shippings(expedition_id);
CREATE INDEX idx_shippings_no_resi ON shippings(no_resi);
