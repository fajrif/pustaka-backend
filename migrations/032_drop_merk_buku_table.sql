-- Migration: drop_merk_buku_table
-- Description: [Add description here]
-- Created: 2025-12-14 09:38:52

-- Ensure UUID extension is enabled (if needed)
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Write your SQL here
DROP TABLE IF EXISTS merk_buku;

-- Drop indexes if exists
DROP INDEX IF EXISTS idx_sales_associates_name;
DROP INDEX IF EXISTS idx_merk_buku_id;
DROP INDEX IF EXISTS idx_merk_buku_kode_merk;
DROP INDEX IF EXISTS idx_merk_buku_nama_merk;
DROP INDEX IF EXISTS idx_merk_buku_user_id;
DROP INDEX IF EXISTS idx_merk_buku_tstamp;

-- Add comments
-- COMMENT ON TABLE table_name IS 'Description';
-- COMMENT ON COLUMN table_name.column_name IS 'Description';
