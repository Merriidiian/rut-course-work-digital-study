CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE IF NOT EXISTS rooms (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    capacity integer NOT NULL CHECK (capacity > 0),
    building_id uuid NOT NULL
);

CREATE TABLE IF NOT EXISTS bookings (
    id uuid PRIMARY KEY,
    room_id uuid NOT NULL REFERENCES rooms(id),
    starts_at timestamptz NOT NULL,
    ends_at timestamptz NOT NULL,
    CHECK (ends_at > starts_at),
    EXCLUDE USING gist (
        room_id WITH =,
        tstzrange(starts_at, ends_at, '[)') WITH &&
    )
);

INSERT INTO rooms (id, name, capacity, building_id) VALUES
('33333333-3333-3333-3333-333333333331','Кабинет 301',25,'55555555-5555-5555-5555-555555555551'),
('33333333-3333-3333-3333-333333333332','Компьютерный класс 405',30,'55555555-5555-5555-5555-555555555552')
ON CONFLICT (id) DO NOTHING;
