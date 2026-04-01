-- UP
ALTER TABLE sales_transactions ADD COLUMN jenis_buku_id UUID REFERENCES jenis_buku(id) ON DELETE SET NULL;
