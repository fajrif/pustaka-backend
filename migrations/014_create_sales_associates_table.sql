-- Migration: Create sales_associates table
-- Description: Creates table to store sales associate data

CREATE TABLE IF NOT EXISTS sales_associates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    address TEXT NOT NULL,
    city_id UUID REFERENCES cities(id) ON DELETE SET NULL ON UPDATE CASCADE,
    area VARCHAR(255),
    phone1 VARCHAR(50) NOT NULL,
    phone2 VARCHAR(50),
    email VARCHAR(255),
    website VARCHAR(255),
    jenis_pembayaran VARCHAR(10) DEFAULT 'T' CHECK (jenis_pembayaran IN ('T', 'K', 'F')),
    join_date DATE NOT NULL,
    end_join_date DATE,
    discount DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_sales_associates_id ON sales_associates(id);
CREATE INDEX IF NOT EXISTS idx_sales_associates_name ON sales_associates(name);
CREATE INDEX IF NOT EXISTS idx_sales_associates_city_id ON sales_associates(city_id);
CREATE INDEX IF NOT EXISTS idx_sales_associates_join_date ON sales_associates(join_date);

-- Add comments to table and columns
COMMENT ON TABLE sales_associates IS 'Table to store sales associate information';
COMMENT ON COLUMN sales_associates.id IS 'UUID primary key (auto-generated)';
COMMENT ON COLUMN sales_associates.name IS 'Name of the sales associate (unique)';
COMMENT ON COLUMN sales_associates.address IS 'Address of the sales associate';
COMMENT ON COLUMN sales_associates.city_id IS 'Reference to cities table';
COMMENT ON COLUMN sales_associates.area IS 'Area/region information';
COMMENT ON COLUMN sales_associates.phone1 IS 'Primary phone number';
COMMENT ON COLUMN sales_associates.phone2 IS 'Secondary phone number';
COMMENT ON COLUMN sales_associates.email IS 'Email address';
COMMENT ON COLUMN sales_associates.website IS 'Website URL';
COMMENT ON COLUMN sales_associates.jenis_pembayaran IS 'Payment type (T/K/F)';
COMMENT ON COLUMN sales_associates.join_date IS 'Date when sales associate joined';
COMMENT ON COLUMN sales_associates.end_join_date IS 'Date when sales associate left (nullable)';
COMMENT ON COLUMN sales_associates.discount IS 'Discount percentage';
COMMENT ON COLUMN sales_associates.created_at IS 'Timestamp of record creation';
COMMENT ON COLUMN sales_associates.updated_at IS 'Timestamp of last update';
