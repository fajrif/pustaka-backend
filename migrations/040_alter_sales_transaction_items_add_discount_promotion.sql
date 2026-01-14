-- UP
-- Migration: Alter sales_transaction_items table - Add discount and promotion fields
-- Description:
--   - promotion: Flat amount deduction from price (e.g., 50 means reduce price by 50)
--   - discount: Percentage discount after promotion (e.g., 10 means 10% off)
--   - Calculation: subtotal = (price - promotion) * (1 - discount/100) * quantity

ALTER TABLE sales_transaction_items ADD COLUMN IF NOT EXISTS promotion NUMERIC(15, 2) NOT NULL DEFAULT 0 CHECK (promotion >= 0);
ALTER TABLE sales_transaction_items ADD COLUMN IF NOT EXISTS discount NUMERIC(5, 2) NOT NULL DEFAULT 0 CHECK (discount >= 0 AND discount <= 100);

-- Add comments
COMMENT ON COLUMN sales_transaction_items.promotion IS 'Flat amount deduction from unit price (default 0)';
COMMENT ON COLUMN sales_transaction_items.discount IS 'Percentage discount applied after promotion (0-100, default 0)';

-- DOWN
-- ALTER TABLE sales_transaction_items DROP COLUMN IF EXISTS promotion;
-- ALTER TABLE sales_transaction_items DROP COLUMN IF EXISTS discount;
