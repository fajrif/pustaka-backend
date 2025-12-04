-- Migration: Create books table
-- Description: Creates table to store book data with multiple foreign keys

CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    jenis_buku_id UUID REFERENCES jenis_buku(id) ON DELETE SET NULL ON UPDATE CASCADE,
    jenjang_studi_id UUID REFERENCES jenjang_studi(id) ON DELETE SET NULL ON UPDATE CASCADE,
    bidang_studi_id UUID REFERENCES bidang_studi(id) ON DELETE SET NULL ON UPDATE CASCADE,
    kelas_id UUID REFERENCES kelas(id) ON DELETE SET NULL ON UPDATE CASCADE,
    publisher_id UUID REFERENCES publishers(id) ON DELETE SET NULL ON UPDATE CASCADE,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_books_id ON books(id);
CREATE INDEX IF NOT EXISTS idx_books_name ON books(name);
CREATE INDEX IF NOT EXISTS idx_books_jenis_buku_id ON books(jenis_buku_id);
CREATE INDEX IF NOT EXISTS idx_books_jenjang_studi_id ON books(jenjang_studi_id);
CREATE INDEX IF NOT EXISTS idx_books_bidang_studi_id ON books(bidang_studi_id);
CREATE INDEX IF NOT EXISTS idx_books_kelas_id ON books(kelas_id);
CREATE INDEX IF NOT EXISTS idx_books_publisher_id ON books(publisher_id);

-- Add comments to table and columns
COMMENT ON TABLE books IS 'Table to store book information';
COMMENT ON COLUMN books.id IS 'UUID primary key (auto-generated)';
COMMENT ON COLUMN books.name IS 'Name of the book';
COMMENT ON COLUMN books.jenis_buku_id IS 'Reference to jenis_buku table (book type)';
COMMENT ON COLUMN books.jenjang_studi_id IS 'Reference to jenjang_studi table (education level)';
COMMENT ON COLUMN books.bidang_studi_id IS 'Reference to bidang_studi table (study field)';
COMMENT ON COLUMN books.kelas_id IS 'Reference to kelas table (class/grade)';
COMMENT ON COLUMN books.publisher_id IS 'Reference to publishers table';
COMMENT ON COLUMN books.price IS 'Price of the book';
COMMENT ON COLUMN books.created_at IS 'Timestamp of record creation';
COMMENT ON COLUMN books.updated_at IS 'Timestamp of last update';
