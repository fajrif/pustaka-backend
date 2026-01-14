-- UP
-- Migration: Alter books table - Add curriculum_id, no_pages, and change kelas from UUID FK to VARCHAR(5)
-- Description:
--   1. Add curriculum_id foreign key to curriculum table
--   2. Add no_pages field (number of pages, default 1)
--   3. Add kelas VARCHAR(5) field to store kelas code directly
--   4. Migrate data from kelas_id to kelas (strip K prefix: K1->1, K2->2, etc.)
--   5. Drop old kelas_id column

-- Step 1: Add new columns
ALTER TABLE books ADD COLUMN IF NOT EXISTS curriculum_id UUID REFERENCES curriculum(id) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE books ADD COLUMN IF NOT EXISTS no_pages INTEGER NOT NULL DEFAULT 1 CHECK (no_pages >= 1);
ALTER TABLE books ADD COLUMN IF NOT EXISTS kelas VARCHAR(5);

-- Step 2: Migrate data from kelas_id to kelas (strip K prefix)
-- K1 -> 1, K2 -> 2, ..., K12 -> 12
-- A, B, ALL stay as-is
UPDATE books b
SET kelas = CASE
    WHEN k.code LIKE 'K%' THEN SUBSTRING(k.code FROM 2)
    ELSE k.code
END
FROM kelas k
WHERE b.kelas_id = k.id;

-- Step 3: Drop the old kelas_id column and its constraints
ALTER TABLE books DROP CONSTRAINT IF EXISTS books_kelas_id_fkey;
DROP INDEX IF EXISTS idx_books_kelas_id;
ALTER TABLE books DROP COLUMN IF EXISTS kelas_id;

-- Step 4: Create new indexes
CREATE INDEX idx_books_curriculum_id ON books(curriculum_id);
CREATE INDEX idx_books_kelas ON books(kelas);
CREATE INDEX idx_books_no_pages ON books(no_pages);

-- Add comments
COMMENT ON COLUMN books.curriculum_id IS 'Reference to curriculum table';
COMMENT ON COLUMN books.no_pages IS 'Number of pages in the book (default 1)';
COMMENT ON COLUMN books.kelas IS 'Class/grade code stored directly (1-12, A, B, ALL)';

-- DOWN
-- Note: This is a complex migration, rollback requires careful handling
-- ALTER TABLE books ADD COLUMN IF NOT EXISTS kelas_id UUID REFERENCES kelas(id) ON DELETE SET NULL ON UPDATE CASCADE;
-- UPDATE books b SET kelas_id = k.id FROM kelas k WHERE b.kelas = k.code OR (b.kelas IS NOT NULL AND k.code = CONCAT('K', b.kelas));
-- ALTER TABLE books DROP COLUMN IF EXISTS kelas;
-- ALTER TABLE books DROP COLUMN IF EXISTS curriculum_id;
-- ALTER TABLE books DROP COLUMN IF EXISTS no_pages;
-- DROP INDEX IF EXISTS idx_books_curriculum_id;
-- DROP INDEX IF EXISTS idx_books_kelas;
-- DROP INDEX IF EXISTS idx_books_no_pages;
