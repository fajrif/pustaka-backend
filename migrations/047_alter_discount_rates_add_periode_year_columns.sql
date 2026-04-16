-- UP
ALTER TABLE discount_rates ADD COLUMN IF NOT EXISTS periode INT NOT NULL DEFAULT 1;
ALTER TABLE discount_rates ADD COLUMN IF NOT EXISTS year VARCHAR(4) NOT NULL DEFAULT EXTRACT(YEAR FROM CURRENT_DATE)::TEXT;
ALTER TABLE discount_rates ADD COLUMN IF NOT EXISTS start_date DATE;
ALTER TABLE discount_rates ADD COLUMN IF NOT EXISTS end_date DATE;

UPDATE discount_rates SET periode = 1, year = EXTRACT(YEAR FROM CURRENT_DATE)::TEXT WHERE name ILIKE '%Early%';
UPDATE discount_rates SET periode = 2, year = EXTRACT(YEAR FROM CURRENT_DATE)::TEXT WHERE name ILIKE '%Secondary%';

CREATE UNIQUE INDEX IF NOT EXISTS idx_discount_rates_periode_year ON discount_rates(periode, year);
