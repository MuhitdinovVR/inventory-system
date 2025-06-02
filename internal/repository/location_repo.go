package repository

import (
	"context"
	"database/sql"
	"errors"
	"inventory-system/internal/models"
)

type LocationRepository struct {
	db *sql.DB
}

func NewLocationRepository(db *sql.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) GetAll(ctx context.Context) ([]models.Location, error) {
	query := `
		SELECT id, address, type
		FROM locations
		ORDER BY address
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var l models.Location
		err := rows.Scan(&l.ID, &l.Address, &l.Type)
		if err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}

	return locations, nil
}

func (r *LocationRepository) GetByID(ctx context.Context, id int) (*models.Location, error) {
	query := `
		SELECT id, address, type
		FROM locations
		WHERE id = $1
	`

	var l models.Location
	err := r.db.QueryRowContext(ctx, query, id).Scan(&l.ID, &l.Address, &l.Type)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &l, nil
}

func (r *LocationRepository) Create(ctx context.Context, location models.Location) (int, error) {
	query := `
		INSERT INTO locations (address, type)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(
		ctx,
		query,
		location.Address,
		location.Type,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *LocationRepository) Update(ctx context.Context, location models.Location) error {
	query := `
		UPDATE locations
		SET address = $1, type = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		location.Address,
		location.Type,
		location.ID,
	)

	return err
}

func (r *LocationRepository) Delete(ctx context.Context, id int) error {
	// Проверяем, есть ли связанные активы
	var hasAssets bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM assets WHERE current_location_id = $1)",
		id,
	).Scan(&hasAssets)

	if err != nil {
		return err
	}

	if hasAssets {
		return errors.New("cannot delete location with assigned assets")
	}

	// Проверяем, есть ли связанные перемещения
	var hasTransfers bool
	err = r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM asset_transfers WHERE from_location_id = $1 OR to_location_id = $1)",
		id,
	).Scan(&hasTransfers)

	if err != nil {
		return err
	}

	if hasTransfers {
		return errors.New("cannot delete location with transfer history")
	}

	_, err = r.db.ExecContext(ctx, "DELETE FROM locations WHERE id = $1", id)
	return err
}

func (r *LocationRepository) GetByType(ctx context.Context, locationType string) ([]models.Location, error) {
	query := `
		SELECT id, address, type
		FROM locations
		WHERE type = $1
		ORDER BY address
	`

	rows, err := r.db.QueryContext(ctx, query, locationType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var l models.Location
		err := rows.Scan(&l.ID, &l.Address, &l.Type)
		if err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}

	return locations, nil
}

func (r *LocationRepository) Exists(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM locations WHERE id = $1)",
		id,
	).Scan(&exists)

	return exists, err
}

func (r *LocationRepository) HasAssets(ctx context.Context, id int) (bool, error) {
	var hasAssets bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM assets WHERE current_location_id = $1)",
		id,
	).Scan(&hasAssets)

	return hasAssets, err
}

func (r *LocationRepository) HasTransfers(ctx context.Context, id int) (bool, error) {
	var hasTransfers bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM asset_transfers WHERE from_location_id = $1 OR to_location_id = $1)",
		id,
	).Scan(&hasTransfers)

	return hasTransfers, err
}
