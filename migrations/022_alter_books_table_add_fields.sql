-- Migration: Alter books table - add description and year fields
-- Description: Adds description and year (mandatory) fields to books table

ALTER TABLE books ADD COLUMN description TEXT;
ALTER TABLE books ADD COLUMN year VARCHAR(4) NOT NULL;

-- Add comments
COMMENT ON COLUMN books.description IS 'Description of the book';
COMMENT ON COLUMN books.year IS 'Year of the book (mandatory)';
