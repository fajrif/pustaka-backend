-- Migration: Alter books table - add author, ISBN, and stock fields
-- Description: Adds author, ISBN, and stock columns to books table (all nullable)

ALTER TABLE books ADD COLUMN author VARCHAR(255);
ALTER TABLE books ADD COLUMN isbn VARCHAR(50);
ALTER TABLE books ADD COLUMN stock INTEGER DEFAULT 0;

-- Create index for ISBN field for faster lookups
CREATE INDEX IF NOT EXISTS idx_books_isbn ON books(isbn);

-- Add comments
COMMENT ON COLUMN books.author IS 'Author of the book (optional)';
COMMENT ON COLUMN books.isbn IS 'ISBN number of the book (optional)';
COMMENT ON COLUMN books.stock IS 'Stock quantity available (default: 0)';
