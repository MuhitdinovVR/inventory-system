package repository

import (
	"context"
	"database/sql"
	"errors"
	"inventory-system/internal/models"
)

type StatusRepository struct {
	db *sql.DB
}

func NewStatusRepository(db *sql.DB) *StatusRepository {
	return &StatusRepository{db: db}
}

func (r *StatusRepository) GetAll(ctx context.Context) ([]models.AssetStatus, error) {
	query := `
		SELECT id, name
		FROM asset_statuses
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []models.AssetStatus
	for rows.Next() {
		var s models.AssetStatus
		err := rows.Scan(&s.ID, &s.Name)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}

	return statuses, nil
}

func (r *StatusRepository) GetByID(ctx context.Context, id int) (*models.AssetStatus, error) {
	query := `
		SELECT id, name
		FROM asset_statuses
		WHERE id = $1
	`

	var s models.AssetStatus
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}

func (r *StatusRepository) GetByName(ctx context.Context, name string) (*models.AssetStatus, error) {
	query := `
		SELECT id, name
		FROM asset_statuses
		WHERE name = $1
	`

	var s models.AssetStatus
	err := r.db.QueryRowContext(ctx, query, name).Scan(&s.ID, &s.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}

func (r *StatusRepository) Exists(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM asset_statuses WHERE id = $1)",
		id,
	).Scan(&exists)

	return exists, err
}
