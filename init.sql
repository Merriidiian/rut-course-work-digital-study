CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE students (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    group_name text NOT NULL,
    email text NOT NULL DEFAULT ''
);

CREATE TABLE teachers (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    group_name text NOT NULL DEFAULT '',
    email text NOT NULL DEFAULT ''
);

CREATE TABLE subjects (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    hours integer NOT NULL CHECK (hours > 0)
);

CREATE TABLE buildings (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    address text NOT NULL
);

CREATE TABLE rooms (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    group_name text NOT NULL DEFAULT '',
    email text NOT NULL DEFAULT '',
    capacity integer NOT NULL CHECK (capacity > 0),
    building_id uuid NOT NULL REFERENCES buildings(id)
);

CREATE TABLE bookings (
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

CREATE TABLE lessons (
    id uuid PRIMARY KEY,
    teacher_id uuid NOT NULL REFERENCES teachers(id),
    subject_id uuid NOT NULL REFERENCES subjects(id),
    room_id uuid NOT NULL REFERENCES rooms(id),
    group_name text NOT NULL,
    starts_at timestamptz NOT NULL,
    ends_at timestamptz NOT NULL,
    status text NOT NULL,
    CHECK (ends_at > starts_at),
    EXCLUDE USING gist (
        teacher_id WITH =,
        tstzrange(starts_at, ends_at, '[)') WITH &&
    ),
    EXCLUDE USING gist (
        group_name WITH =,
        tstzrange(starts_at, ends_at, '[)') WITH &&
    )
);

INSERT INTO buildings (id, name, address) VALUES
    ('55555555-5555-5555-5555-555555555551', 'Учебный корпус 1', 'Москва, ул. Образцова, 9'),
    ('55555555-5555-5555-5555-555555555552', 'Учебный корпус 2', 'Москва, ул. Образцова, 9, строение 2');

INSERT INTO students (id, name, group_name, email) VALUES
('11111111-1111-1111-1111-111111111111','Иванов Иван','ИТ-271','ivanov@university.ru'),
('11111111-1111-1111-1111-111111111112','Петрова Анна','ИТ-271','petrova@university.ru'),
('11111111-1111-1111-1111-111111111113','Смирнов Алексей','ИТ-272','smirnov@university.ru');
INSERT INTO teachers (id, name, group_name, email) VALUES
('22222222-2222-2222-2222-222222222221','Соколов Дмитрий','','sokolov@university.ru'),
('22222222-2222-2222-2222-222222222222','Волкова Елена','','volkova@university.ru');
INSERT INTO rooms (id, name, group_name, email, capacity, building_id) VALUES
('33333333-3333-3333-3333-333333333331','Кабинет 301','','',25,'55555555-5555-5555-5555-555555555551'),
('33333333-3333-3333-3333-333333333332','Компьютерный класс 405','','',30,'55555555-5555-5555-5555-555555555552');
INSERT INTO subjects VALUES
('44444444-4444-4444-4444-444444444441','Проектирование корпоративных приложений',108),
('44444444-4444-4444-4444-444444444442','Базы данных',144),
('44444444-4444-4444-4444-444444444443','Алгоритмы и структуры данных',144),
('44444444-4444-4444-4444-444444444444','Компьютерные сети',108),
('44444444-4444-4444-4444-444444444445','Разработка веб-приложений',108),
('44444444-4444-4444-4444-444444444446','Информационная безопасность',72);
