CREATE TABLE IF NOT EXISTS students (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    group_name text NOT NULL,
    email text NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS teachers (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    email text NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS subjects (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    hours integer NOT NULL CHECK (hours > 0)
);

INSERT INTO students (id, name, group_name, email) VALUES
('11111111-1111-1111-1111-111111111111','Иванов Иван','ИТ-271','ivanov@university.ru'),
('11111111-1111-1111-1111-111111111112','Петрова Анна','ИТ-271','petrova@university.ru'),
('11111111-1111-1111-1111-111111111113','Смирнов Алексей','ИТ-272','smirnov@university.ru')
ON CONFLICT (id) DO NOTHING;

INSERT INTO teachers (id, name, email) VALUES
('22222222-2222-2222-2222-222222222221','Соколов Дмитрий','sokolov@university.ru'),
('22222222-2222-2222-2222-222222222222','Волкова Елена','volkova@university.ru')
ON CONFLICT (id) DO NOTHING;

INSERT INTO subjects VALUES
('44444444-4444-4444-4444-444444444441','Проектирование корпоративных приложений',108),
('44444444-4444-4444-4444-444444444442','Базы данных',144),
('44444444-4444-4444-4444-444444444443','Алгоритмы и структуры данных',144),
('44444444-4444-4444-4444-444444444444','Компьютерные сети',108),
('44444444-4444-4444-4444-444444444445','Разработка веб-приложений',108),
('44444444-4444-4444-4444-444444444446','Информационная безопасность',72)
ON CONFLICT (id) DO NOTHING;
