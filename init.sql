CREATE EXTENSION IF NOT EXISTS btree_gist;
CREATE TABLE students(id uuid PRIMARY KEY DEFAULT gen_random_uuid(),name text NOT NULL,group_name text NOT NULL DEFAULT '',email text NOT NULL DEFAULT '',capacity integer NOT NULL DEFAULT 0);
CREATE TABLE teachers(LIKE students INCLUDING ALL);
CREATE TABLE rooms(LIKE students INCLUDING ALL);
ALTER TABLE rooms ADD CHECK (capacity>0);
CREATE TABLE subjects(id uuid PRIMARY KEY DEFAULT gen_random_uuid(),name text NOT NULL,hours integer NOT NULL CHECK(hours>0));
CREATE TABLE bookings(id uuid PRIMARY KEY,room_id uuid NOT NULL REFERENCES rooms(id),starts_at timestamptz NOT NULL,ends_at timestamptz NOT NULL,CHECK(ends_at>starts_at),EXCLUDE USING gist(room_id WITH =,tstzrange(starts_at,ends_at,'[)') WITH &&));
CREATE TABLE lessons(id uuid PRIMARY KEY,teacher_id uuid NOT NULL REFERENCES teachers(id),subject_id uuid NOT NULL REFERENCES subjects(id),room_id uuid NOT NULL REFERENCES rooms(id),group_name text NOT NULL,starts_at timestamptz NOT NULL,ends_at timestamptz NOT NULL,status text NOT NULL,CHECK(ends_at>starts_at),EXCLUDE USING gist(teacher_id WITH =,tstzrange(starts_at,ends_at,'[)') WITH &&),EXCLUDE USING gist(group_name WITH =,tstzrange(starts_at,ends_at,'[)') WITH &&));
INSERT INTO students VALUES
('11111111-1111-1111-1111-111111111111','Иванов Иван','ИТ-271','ivanov@university.ru',0),
('11111111-1111-1111-1111-111111111112','Петрова Анна','ИТ-271','petrova@university.ru',0),
('11111111-1111-1111-1111-111111111113','Смирнов Алексей','ИТ-272','smirnov@university.ru',0);
INSERT INTO teachers VALUES
('22222222-2222-2222-2222-222222222221','Соколов Дмитрий','','sokolov@university.ru',0),
('22222222-2222-2222-2222-222222222222','Волкова Елена','','volkova@university.ru',0);
INSERT INTO rooms VALUES
('33333333-3333-3333-3333-333333333331','Кабинет 301','','',25),
('33333333-3333-3333-3333-333333333332','Компьютерный класс 405','','',30);
INSERT INTO subjects VALUES
('44444444-4444-4444-4444-444444444441','Проектирование корпоративных приложений',108),
('44444444-4444-4444-4444-444444444442','Базы данных',144),
('44444444-4444-4444-4444-444444444443','Алгоритмы и структуры данных',144),
('44444444-4444-4444-4444-444444444444','Компьютерные сети',108),
('44444444-4444-4444-4444-444444444445','Разработка веб-приложений',108),
('44444444-4444-4444-4444-444444444446','Информационная безопасность',72);
