package buildings

import (
	"context"
	"database/sql"
)

type Repository struct {
	database *sql.DB
}

func (repository *Repository) List(ctx context.Context) ([]Building, error) {
	rows, err := repository.database.QueryContext(ctx,
		"SELECT id, name, address FROM buildings ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []Building{}
	for rows.Next() {
		var entity Building
		if err := rows.Scan(&entity.ID, &entity.Name, &entity.Address); err != nil {
			return nil, err
		}
		result = append(result, entity)
	}
	return result, rows.Err()
}

func (repository *Repository) Find(ctx context.Context, id string) (Building, error) {
	var entity Building
	err := repository.database.QueryRowContext(ctx,
		"SELECT id, name, address FROM buildings WHERE id = $1", id).Scan(&entity.ID, &entity.Name, &entity.Address)
	return entity, err
}

func (repository *Repository) Save(ctx context.Context, entity Building) (Building, error) {
	var err error
	if entity.ID == "" {
		err = repository.database.QueryRowContext(ctx,
			"INSERT INTO buildings (name, address) VALUES ($1, $2) RETURNING id",
			entity.Name, entity.Address).Scan(&entity.ID)
	} else {
		err = repository.database.QueryRowContext(ctx,
			"UPDATE buildings SET name = $2, address = $3 WHERE id = $1 RETURNING id",
			entity.ID, entity.Name, entity.Address).Scan(&entity.ID)
	}
	return entity, err
}

func (repository *Repository) Delete(ctx context.Context, id string) (bool, error) {
	result, err := repository.database.ExecContext(ctx, "DELETE FROM buildings WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count > 0, err
}
