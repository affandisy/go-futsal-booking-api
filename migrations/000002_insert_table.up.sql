-- 000002_insert_table.up.sql

-- Insert into roles
INSERT INTO roles (id, role_name)
VALUES
    (1, 'ADMIN'),
    (2, 'CUSTOMER');


-- Insert into users
INSERT INTO users (id, role_id, full_name, email, is_verified, password, age, address)
VALUES
    (1, 1, 'Admin Utama', 'admin@example.com', true, 'admin123', 19, 'Jl. Gatot Subroto No. 1, Indramayu'),
    (2, 2, 'Jane Doe', 'janedoe@example.com', true, 'password123', 20, 'Jl. Ismail No. 5, Indramayu'),
    (3, 2, 'Shima Rin', 'shimarin@example.com', true, 'password123', 18, 'Jl. Sudirman No. 6, Indramayu');


-- Insert into venues
INSERT INTO venues (id, name, address, city)
VALUES
    (1, 'VIP FUTSAL', 'Jl. Olahraga No. 10', 'Indramayu'),
    (2, 'Walang Futsal', 'Jl. Geulis No. 20', 'Indramayu'),
    (3, 'Sunday Futsal', 'Jl. Kitabisa No. 14', 'Indramayu');

-- Insert into fields
INSERT INTO fields (id, venue_id, name, type)
VALUES
    (1, 1, 'Lapangan A', 'SINTETIS'),
    (2, 1, 'Lapangan B', 'VINYL'),
    (3, 1, 'Lapangan C', 'BETON'),
    (4, 2, 'Lapangan Walang A', 'SINTETIS'),
    (5, 2, 'Lapangan Walang B', 'VINYL'),
    (6, 3, 'Lapangan Sunday A', 'SINTETIS');

-- Insert into schedules
INSERT INTO schedules (id, field_id, day_of_week, start_time, end_time, price)
VALUES 
    -- VIP FUTSAL
    (1, 1, 6, '19:00:00', '20:00:00', 150000.00), -- Sabtu, 19.00
    (2, 1, 7, '19:00:00', '20:00:00', 200000.00), -- Minggu, 19.00
    -- WALANG FUTSAL
    (3, 2, 6, '19:00:00', '20:00:00', 120000.00), -- Sabtu, 19.00
    (4, 2, 7, '19:00:00', '20:00:00', 130000.00), -- Minggu, 19.00
    -- SUNDAY FUTSAL
    (5, 3, 6, '19:00:00', '20:00:00', 100000.00), -- Sabtu, 19.00
    (6, 3, 7, '19:00:00', '20:00:00', 120000.00); -- Minggu, 19.00

-- Insert into bookings
INSERT INTO bookings (id, user_id, schedule_id, booking_date, status, total_price)
VALUES
    (1, 2, 1, '2025-11-10', 'PENDING', 150000.00),
    (2, 3, 3, '2025-11-11', 'PENDING', 120000.00);

-- Insert into payments
INSERT INTO payments (id, booking_id, payment_method, amount, status)
VALUES
    (1, 1, 'E_WALLET', 150000.00, 'SUCCESS'),
    (2, 2, 'CASH', 120000.00, 'SUCCESS');

SELECT setval(pg_get_serial_sequence('users', 'id'), (SELECT MAX(id) FROM users));
SELECT setval(pg_get_serial_sequence('venues', 'id'), (SELECT MAX(id) FROM venues));
SELECT setval(pg_get_serial_sequence('fields', 'id'), (SELECT MAX(id) FROM fields));
SELECT setval(pg_get_serial_sequence('schedules', 'id'), (SELECT MAX(id) FROM schedules));
SELECT setval(pg_get_serial_sequence('bookings', 'id'), (SELECT MAX(id) FROM bookings));
