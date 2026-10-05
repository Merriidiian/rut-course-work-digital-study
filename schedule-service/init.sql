CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE IF NOT EXISTS lessons (
    id uuid PRIMARY KEY,
    teacher_id uuid NOT NULL,
    subject_id uuid NOT NULL,
    room_id uuid NOT NULL,
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
