BEGIN;

CREATE TABLE IF NOT EXISTS buildings (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    address text NOT NULL
);

INSERT INTO buildings (id, name, address) VALUES
    ('55555555-5555-5555-5555-555555555551', 'Учебный корпус 1', 'Москва, ул. Образцова, 9'),
    ('55555555-5555-5555-5555-555555555552', 'Учебный корпус 2', 'Москва, ул. Образцова, 9, строение 2')
ON CONFLICT (id) DO NOTHING;

ALTER TABLE rooms ADD COLUMN IF NOT EXISTS building_id uuid REFERENCES buildings(id);

UPDATE rooms
SET building_id = '55555555-5555-5555-5555-555555555551'
WHERE building_id IS NULL;

UPDATE rooms
SET building_id = '55555555-5555-5555-5555-555555555552'
WHERE id = '33333333-3333-3333-3333-333333333332';

ALTER TABLE rooms ALTER COLUMN building_id SET NOT NULL;

COMMIT;
