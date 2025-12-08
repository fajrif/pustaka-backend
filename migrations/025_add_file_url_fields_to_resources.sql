-- Migration: Add file URL fields to resources
-- Description: Adds photo_url, image_url, file_url, and logo_url fields to support file uploads

-- UP

-- Add photo_url to users table
ALTER TABLE users
ADD COLUMN photo_url VARCHAR(500);

COMMENT ON COLUMN users.photo_url IS 'URL path to user profile photo';

-- Add image_url and file_url to books table
ALTER TABLE books
ADD COLUMN image_url VARCHAR(500),
ADD COLUMN file_url VARCHAR(500);

COMMENT ON COLUMN books.image_url IS 'URL path to book cover image';
COMMENT ON COLUMN books.file_url IS 'URL path to book file (PDF, etc.)';

-- Add logo_url and file_url to publishers table
ALTER TABLE publishers
ADD COLUMN logo_url VARCHAR(500),
ADD COLUMN file_url VARCHAR(500);

COMMENT ON COLUMN publishers.logo_url IS 'URL path to publisher logo';
COMMENT ON COLUMN publishers.file_url IS 'URL path to publisher related files';

-- Add logo_url and file_url to expeditions table
ALTER TABLE expeditions
ADD COLUMN logo_url VARCHAR(500),
ADD COLUMN file_url VARCHAR(500);

COMMENT ON COLUMN expeditions.logo_url IS 'URL path to expedition logo';
COMMENT ON COLUMN expeditions.file_url IS 'URL path to expedition related files';

-- Add photo_url and file_url to sales_associates table
ALTER TABLE sales_associates
ADD COLUMN photo_url VARCHAR(500),
ADD COLUMN file_url VARCHAR(500);

COMMENT ON COLUMN sales_associates.photo_url IS 'URL path to sales associate photo';
COMMENT ON COLUMN sales_associates.file_url IS 'URL path to sales associate related files';
