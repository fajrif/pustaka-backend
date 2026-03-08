-- UP
ALTER TABLE sales_transactions
    ADD COLUMN IF NOT EXISTS periode INT NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS year VARCHAR(4) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS curriculum_id UUID REFERENCES curriculum(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS merk_buku_id UUID REFERENCES merk_buku(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS jenjang_studi_id UUID REFERENCES jenjang_studi(id) ON DELETE SET NULL;

-- DOWN
-- ALTER TABLE sales_transactions
--     DROP COLUMN IF EXISTS periode,
--     DROP COLUMN IF EXISTS year,
--     DROP COLUMN IF EXISTS curriculum_id,
--     DROP COLUMN IF EXISTS merk_buku_id,
--     DROP COLUMN IF EXISTS jenjang_studi_id;
