-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Projects table
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    no_sp2k VARCHAR(255) NOT NULL,
    no_perjanjian VARCHAR(255),
    no_amandemen VARCHAR(255),
    tanggal_perjanjian DATE,
    judul_pekerjaan VARCHAR(500) NOT NULL,
    jangka_waktu INTEGER,
    tanggal_mulai DATE NOT NULL,
    tanggal_selesai DATE,
    nilai_pekerjaan NUMERIC(20, 2) NOT NULL,
    management_fee NUMERIC(20, 2),
    tarif_management_fee_persen NUMERIC(5, 2),
    client VARCHAR(255),
    pic_client VARCHAR(255),
    contact_client VARCHAR(255),
    alamat_client TEXT,
    jenis_kontrak VARCHAR(100),
    status_kontrak VARCHAR(50) DEFAULT 'Active',
    created_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Cost Types table
CREATE TABLE cost_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama_biaya VARCHAR(255) NOT NULL,
    kode VARCHAR(50),
    deskripsi TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Budget Items table
CREATE TABLE budget_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    no_sp2k VARCHAR(255),
    cost_type_id UUID REFERENCES cost_types(id),
    jenis_biaya_name VARCHAR(255),
    kategori_anggaran VARCHAR(50),
    total_anggaran NUMERIC(20, 2),
    deskripsi_anggaran TEXT,
    periode_bulan VARCHAR(7),
    jumlah_anggaran NUMERIC(20, 2),
    bulan_ke INTEGER,
    parent_budget_id UUID REFERENCES budget_items(id) ON DELETE CASCADE,
    is_parent BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Transactions table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    no_sp2k VARCHAR(255),
    tanggal_transaksi DATE NOT NULL,
    tanggal_po_tagihan DATE,
    bulan_realisasi VARCHAR(7),
    cost_type_id UUID REFERENCES cost_types(id),
    jenis_biaya_name VARCHAR(255),
    deskripsi_realisasi TEXT,
    jumlah_realisasi NUMERIC(20, 2) NOT NULL,
    persentase_management_fee NUMERIC(5, 2),
    nilai_management_fee NUMERIC(20, 2),
    jumlah_tenaga_kerja INTEGER,
    bukti_transaksi_url TEXT,
    created_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX idx_projects_status ON projects(status_kontrak);
CREATE INDEX idx_budget_items_project ON budget_items(project_id);
CREATE INDEX idx_budget_items_parent ON budget_items(parent_budget_id);
CREATE INDEX idx_transactions_project ON transactions(project_id);
CREATE INDEX idx_transactions_month ON transactions(bulan_realisasi);

-- Insert default admin user (password: admin123)
INSERT INTO users (email, password_hash, full_name, role)
VALUES ('admin@budgetwise.com', '$2a$10$b/ipDn9ncPDMmbrf2J2rZOE0IA8Rv4QuTyDHgEOd1gOYkV7cwlQaS', 'Administrator', 'admin')
