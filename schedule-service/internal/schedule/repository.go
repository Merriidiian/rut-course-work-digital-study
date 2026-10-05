package schedule

import (
	"context"
	"database/sql"
)

type Repository struct {
	database *sql.DB
}

func (repository *Repository) NewID(ctx context.Context) (string, error) {
	var id string
	err := repository.database.QueryRowContext(ctx, "SELECT gen_random_uuid()").Scan(&id)
	return id, err
}

func (repository *Repository) List(ctx context.Context, group, teacher, subject string) ([]Lesson, error) {
	rows, err := repository.database.QueryContext(ctx, `
        SELECT id, teacher_id, subject_id, room_id, group_name, starts_at, ends_at, status
        FROM lessons
        WHERE ($1 = '' OR group_name = $1)
          AND ($2 = '' OR teacher_id::text = $2)
          AND ($3 = '' OR subject_id::text = $3)
        ORDER BY starts_at`, group, teacher, subject)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lessons := []Lesson{}
	for rows.Next() {
		var lesson Lesson
		err := rows.Scan(&lesson.ID, &lesson.TeacherID, &lesson.SubjectID, &lesson.RoomID,
			&lesson.Group, &lesson.StartsAt, &lesson.EndsAt, &lesson.Status)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}
	return lessons, rows.Err()
}

func (repository *Repository) Save(ctx context.Context, lesson Lesson) error {
	_, err := repository.database.ExecContext(ctx, `
        INSERT INTO lessons
            (id, teacher_id, subject_id, room_id, group_name, starts_at, ends_at, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		lesson.ID, lesson.TeacherID, lesson.SubjectID, lesson.RoomID,
		lesson.Group, lesson.StartsAt, lesson.EndsAt, lesson.Status)
	return err
}

func (repository *Repository) Delete(ctx context.Context, id string) (bool, error) {
	result, err := repository.database.ExecContext(ctx, "DELETE FROM lessons WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count > 0, err
}
