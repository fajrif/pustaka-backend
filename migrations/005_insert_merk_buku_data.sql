-- Migration: Insert merk_buku data (Version 2 - with UUID user_id)
-- Description: Inserts initial data for book brands with UUID foreign key to users

-- Note: This assumes users table has been created and populated
-- User UUIDs:
-- KURNIA: a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11
-- ALBERT: b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22
-- NOVI:   c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33
-- DEVITA: d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44
-- FENNY:  e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55

-- Insert data into merk_buku table
INSERT INTO merk_buku (kode_merk, nama_merk, bantuan_promosi, user_id, tstamp) VALUES
('GHD', 'GRAHADI', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2006-12-15 10:39:06.217'),
('MCL', 'MENTARI COVER LAMA', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2006-01-23 13:34:06.327'),
('MDR', 'MENTARI KB', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2013-12-27 17:13:02.203'),
('MEN', 'MENTARI', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2006-01-23 13:33:52.703'),
('MPR', 'WAJAR KB  2013', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2014-01-06 11:16:30.827'),
('MXX', 'MAXXI', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2010-06-14 21:54:22.467'),
('PRO', 'PROYEKSI ', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2006-07-22 13:52:54.890'),
('SCL', 'SISWA COVER LAMA', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2006-06-28 08:06:02.437'),
('SIS', 'SISWA', 0, 'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', '2005-01-03 21:38:22.203'),
('SKR', 'SEKAR', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2006-06-29 08:45:29.843'),
('SMT', 'SMART MEDIA', 0, 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', '2005-10-17 11:21:21.483'),
('TCL', 'TUNTAS COVER LAMA', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2006-05-11 08:32:59.327'),
('TUN', 'TUNTAS', 0, 'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', '2005-01-03 21:38:44.153'),
('UAN', 'SIAP UJIAN NASIONAL', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2008-10-27 11:37:52.780'),
('WCL', 'WAJAR', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2015-08-18 09:00:55.843'),
('WJR', 'WAJAR', 0, 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', '2006-02-03 09:31:17.733'),
('MUS', 'KARTIKA PRIMA', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2002-01-01 00:55:16.920'),
('KB', 'TUNTAS KB 2013', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2013-12-30 08:55:33.577'),
('WKB', 'WAJAR KURIKULUM 2013', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2014-04-07 13:37:08.767'),
('KPB', 'KARTIKA PRIMA KB ', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2014-04-29 11:19:53.767'),
('KBR', 'K.BULAN RAMADHAN', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2014-05-26 15:10:50.203'),
('FT', 'FATTAH', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2017-04-27 15:34:30.750'),
('HK', 'HIKMAH', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2023-12-30 08:34:22.093'),
('KEA', 'SMK KEAHLIAN', 0, 'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', '2024-06-03 13:21:43.357'),
('PM', 'PRIMA', 0, 'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', '2024-08-26 13:13:35.267'),
('SAJ', 'UAN KUMER', 0, 'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', '2024-10-15 09:16:01.810'),
('JUR', 'JURNAL', 0, 'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', '2025-05-08 15:25:24.040'),
('TKA', 'TES KEMAMPUAN DASAR', 0, 'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', '2025-09-03 10:34:58.430'),
('KAR', 'KARTIKA PRATAMA', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2002-01-01 00:55:30.360'),
('PRA', 'PRAKARYA', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2014-04-08 08:15:57.687'),
('KK', 'KOMPETENSI KEAHLIAN', 0, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2015-06-10 10:37:30.960'),
('FJR', 'FAJAR', 0, 'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', '2024-01-09 11:33:00.340')
ON CONFLICT (kode_merk) DO NOTHING;

-- Verify the insert with user information
SELECT 
    mb.kode_merk,
    mb.nama_merk,
    u.full_name as created_by,
    u.email,
    mb.tstamp
FROM merk_buku mb
LEFT JOIN users u ON mb.user_id = u.id
ORDER BY mb.kode_merk;

-- Count records
SELECT COUNT(*) as total_records FROM merk_buku;
