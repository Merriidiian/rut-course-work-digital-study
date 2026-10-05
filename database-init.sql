DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'core_user') THEN
        CREATE ROLE core_user LOGIN PASSWORD 'core';
    END IF;
END
$$;

SELECT 'CREATE DATABASE university_core OWNER core_user'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'university_core')
\gexec

REVOKE CONNECT ON DATABASE university_core FROM PUBLIC;
GRANT CONNECT ON DATABASE university_core TO core_user;

DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'rooms_user') THEN
        CREATE ROLE rooms_user LOGIN PASSWORD 'rooms';
    END IF;
END
$$;

SELECT 'CREATE DATABASE university_rooms OWNER rooms_user'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'university_rooms')
\gexec

REVOKE CONNECT ON DATABASE university_rooms FROM PUBLIC;
GRANT CONNECT ON DATABASE university_rooms TO rooms_user;

DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'buildings_user') THEN
        CREATE ROLE buildings_user LOGIN PASSWORD 'buildings';
    END IF;
END
$$;

SELECT 'CREATE DATABASE university_buildings OWNER buildings_user'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'university_buildings')
\gexec

REVOKE CONNECT ON DATABASE university_buildings FROM PUBLIC;
GRANT CONNECT ON DATABASE university_buildings TO buildings_user;

DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'schedule_user') THEN
        CREATE ROLE schedule_user LOGIN PASSWORD 'schedule';
    END IF;
END
$$;

SELECT 'CREATE DATABASE university_schedule OWNER schedule_user'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'university_schedule')
\gexec

REVOKE CONNECT ON DATABASE university_schedule FROM PUBLIC;
GRANT CONNECT ON DATABASE university_schedule TO schedule_user;
