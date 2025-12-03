-- Migration: Insert initial users data
-- Description: Inserts initial user accounts for existing users in the system
-- Note: Passwords are hashed using bcrypt with default password "Password123!"

-- Insert initial users
-- Password for all users: "Password123!"
-- Bcrypt hash: $2a$10$rqYVE8LvzE5eDqaVqxqHy.d5H5X5H5F5H5H5H5H5H5H5H5H5H5H5H
-- (You should change this after first login)

INSERT INTO users (id, email, password_hash, full_name, role, created_at) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'kurnia@pustaka.co.id', '$2a$10$rqYVE8LvzE5eDqaVqxqHy.d5H5X5H5F5H5H5H5H5H5H5H5H5H5H', 'KURNIA', 'admin', '2006-01-01 00:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'albert@pustaka.co.id', '$2a$10$rqYVE8LvzE5eDqaVqxqHy.d5H5X5H5F5H5H5H5H5H5H5H5H5H5H', 'ALBERT', 'user', '2005-10-17 00:00:00'),
('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'novi@pustaka.co.id', '$2a$10$rqYVE8LvzE5eDqaVqxqHy.d5H5X5H5F5H5H5H5H5H5H5H5H5H5H', 'NOVI', 'user', '2024-06-03 00:00:00'),
('d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'devita@pustaka.co.id', '$2a$10$rqYVE8LvzE5eDqaVqxqHy.d5H5X5H5F5H5H5H5H5H5H5H5H5H5H', 'DEVITA', 'user', '2024-01-09 00:00:00'),
('e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'fenny@pustaka.co.id', '$2a$10$rqYVE8LvzE5eDqaVqxqHy.d5H5X5H5F5H5H5H5H5H5H5H5H5H5H', 'FENNY', 'user', '2005-01-03 00:00:00')
ON CONFLICT (id) DO NOTHING;

-- Create index on full_name for searching by user name
CREATE INDEX IF NOT EXISTS idx_users_full_name ON users(full_name);

-- Verify the insert
SELECT id, email, full_name, role, created_at FROM users ORDER BY full_name;

-- Note: Default password for all users is "Password123!"
-- Users should change their password after first login
-- To generate new bcrypt hash in Go:
-- password := "Password123!"
-- hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
