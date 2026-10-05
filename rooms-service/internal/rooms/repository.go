package rooms

import (
	"context"
	"database/sql"
)

type Repository struct {
	database *sql.DB
}

func (repository *Repository) List(ctx context.Context, buildingID string) ([]Room, error) {
	rows, err := repository.database.QueryContext(ctx,
		"SELECT id, name, building_id, capacity FROM rooms WHERE ($1 = '' OR building_id::text = $1) ORDER BY name", buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []Room{}
	for rows.Next() {
		var entity Room
		if err := rows.Scan(&entity.ID, &entity.Name, &entity.BuildingID, &entity.Capacity); err != nil {
			return nil, err
		}
		result = append(result, entity)
	}
	return result, rows.Err()
}

func (repository *Repository) Find(ctx context.Context, id string) (Room, error) {
	var entity Room
	err := repository.database.QueryRowContext(ctx,
		"SELECT id, name, building_id, capacity FROM rooms WHERE id = $1", id).Scan(&entity.ID, &entity.Name, &entity.BuildingID, &entity.Capacity)
	return entity, err
}

func (repository *Repository) Save(ctx context.Context, entity Room) (Room, error) {
	var err error
	if entity.ID == "" {
		err = repository.database.QueryRowContext(ctx,
			"INSERT INTO rooms (name, building_id, capacity) VALUES ($1, $2, $3) RETURNING id",
			entity.Name, entity.BuildingID, entity.Capacity).Scan(&entity.ID)
	} else {
		err = repository.database.QueryRowContext(ctx,
			"UPDATE rooms SET name = $2, building_id = $3, capacity = $4 WHERE id = $1 RETURNING id",
			entity.ID, entity.Name, entity.BuildingID, entity.Capacity).Scan(&entity.ID)
	}
	return entity, err
}

func (repository *Repository) Delete(ctx context.Context, id string) (bool, error) {
	result, err := repository.database.ExecContext(ctx, "DELETE FROM rooms WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count > 0, err
}
