-- Clients table
CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    contact_name VARCHAR(255),
    phone VARCHAR(50),
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Contract Types table
CREATE TABLE contract_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    code VARCHAR(50),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- alter Projects table
ALTER TABLE projects DROP COLUMN client;
ALTER TABLE projects DROP COLUMN pic_client;
ALTER TABLE projects DROP COLUMN contact_client;
ALTER TABLE projects DROP COLUMN alamat_client;
ALTER TABLE projects DROP COLUMN jenis_kontrak;

ALTER TABLE projects ADD COLUMN client_id UUID REFERENCES clients(id) ON DELETE SET NULL;
ALTER TABLE projects ADD COLUMN contract_type_id UUID REFERENCES contract_types(id) ON DELETE SET NULL;

-- Indexes for better performance
CREATE INDEX idx_projects_client ON projects(client_id);
CREATE INDEX idx_projects_contract_type ON projects(contract_type_id);
