CREATE TABLE IF NOT EXISTS buildings (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    address text NOT NULL
);

INSERT INTO buildings (id, name, address) VALUES
    ('55555555-5555-5555-5555-555555555551', 'Учебный корпус 1', 'Москва, ул. Образцова, 9'),
    ('55555555-5555-5555-5555-555555555552', 'Учебный корпус 2', 'Москва, ул. Образцова, 9, строение 2')
ON CONFLICT (id) DO NOTHING;
