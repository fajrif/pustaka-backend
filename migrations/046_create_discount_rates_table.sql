-- UP
CREATE TABLE discount_rates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    discount DECIMAL(5,2) NOT NULL DEFAULT 0,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO discount_rates (name, discount, description) VALUES
    ('Early Payment Discount', 8.00, 'Discount for payments made on or before the due date'),
    ('Secondary Payment Discount', 5.00, 'Discount for payments made on or before the secondary due date');
