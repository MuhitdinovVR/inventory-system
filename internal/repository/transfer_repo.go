package repository

import (
	"context"
	"database/sql"
	"errors"
	"inventory-system/internal/models"
	"time"
)

type TransferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{db: db}
}

func (r *TransferRepository) GetAll(ctx context.Context) ([]models.AssetTransfer, error) {
	query := `
		SELECT t.id, t.asset_id, a.name as asset_name,
		       t.employee_id, e.full_name as employee_name,
		       t.from_location_id, fl.address as from_location,
		       t.to_location_id, tl.address as to_location,
		       t.transfer_date, t.notes
		FROM asset_transfers t
		JOIN assets a ON t.asset_id = a.id
		JOIN employees e ON t.employee_id = e.id
		JOIN locations fl ON t.from_location_id = fl.id
		JOIN locations tl ON t.to_location_id = tl.id
		ORDER BY t.transfer_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []models.AssetTransfer
	for rows.Next() {
		var t models.AssetTransfer
		err := rows.Scan(
			&t.ID,
			&t.AssetID,
			&t.AssetName,
			&t.EmployeeID,
			&t.EmployeeName,
			&t.FromLocationID,
			&t.FromLocation,
			&t.ToLocationID,
			&t.ToLocation,
			&t.TransferDate,
			&t.Notes,
		)
		if err != nil {
			return nil, err
		}

		transfers = append(transfers, t)
	}

	return transfers, nil
}

func (r *TransferRepository) GetByID(ctx context.Context, id int) (*models.AssetTransfer, error) {
	query := `
		SELECT t.id, t.asset_id, a.name as asset_name,
		       t.employee_id, e.full_name as employee_name,
		       t.from_location_id, fl.address as from_location,
		       t.to_location_id, tl.address as to_location,
		       t.transfer_date, t.notes
		FROM asset_transfers t
		JOIN assets a ON t.asset_id = a.id
		JOIN employees e ON t.employee_id = e.id
		JOIN locations fl ON t.from_location_id = fl.id
		JOIN locations tl ON t.to_location_id = tl.id
		WHERE t.id = $1
	`

	var t models.AssetTransfer
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.AssetID,
		&t.AssetName,
		&t.EmployeeID,
		&t.EmployeeName,
		&t.FromLocationID,
		&t.FromLocation,
		&t.ToLocationID,
		&t.ToLocation,
		&t.TransferDate,
		&t.Notes,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}

func (r *TransferRepository) Create(ctx context.Context, transfer models.AssetTransfer) (int, error) {
	query := `
		INSERT INTO asset_transfers (
			asset_id, employee_id, from_location_id, 
			to_location_id, transfer_date, notes
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(
		ctx,
		query,
		transfer.AssetID,
		transfer.EmployeeID,
		transfer.FromLocationID,
		transfer.ToLocationID,
		transfer.TransferDate,
		transfer.Notes,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *TransferRepository) GetByAsset(ctx context.Context, assetID int) ([]models.AssetTransfer, error) {
	query := `
		SELECT t.id, t.asset_id, a.name as asset_name,
		       t.employee_id, e.full_name as employee_name,
		       t.from_location_id, fl.address as from_location,
		       t.to_location_id, tl.address as to_location,
		       t.transfer_date, t.notes
		FROM asset_transfers t
		JOIN assets a ON t.asset_id = a.id
		JOIN employees e ON t.employee_id = e.id
		JOIN locations fl ON t.from_location_id = fl.id
		JOIN locations tl ON t.to_location_id = tl.id
		WHERE t.asset_id = $1
		ORDER BY t.transfer_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []models.AssetTransfer
	for rows.Next() {
		var t models.AssetTransfer
		err := rows.Scan(
			&t.ID,
			&t.AssetID,
			&t.AssetName,
			&t.EmployeeID,
			&t.EmployeeName,
			&t.FromLocationID,
			&t.FromLocation,
			&t.ToLocationID,
			&t.ToLocation,
			&t.TransferDate,
			&t.Notes,
		)
		if err != nil {
			return nil, err
		}

		transfers = append(transfers, t)
	}

	return transfers, nil
}

func (r *TransferRepository) GetByEmployee(ctx context.Context, employeeID int) ([]models.AssetTransfer, error) {
	query := `
		SELECT t.id, t.asset_id, a.name as asset_name,
		       t.employee_id, e.full_name as employee_name,
		       t.from_location_id, fl.address as from_location,
		       t.to_location_id, tl.address as to_location,
		       t.transfer_date, t.notes
		FROM asset_transfers t
		JOIN assets a ON t.asset_id = a.id
		JOIN employees e ON t.employee_id = e.id
		JOIN locations fl ON t.from_location_id = fl.id
		JOIN locations tl ON t.to_location_id = tl.id
		WHERE t.employee_id = $1
		ORDER BY t.transfer_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []models.AssetTransfer
	for rows.Next() {
		var t models.AssetTransfer
		err := rows.Scan(
			&t.ID,
			&t.AssetID,
			&t.AssetName,
			&t.EmployeeID,
			&t.EmployeeName,
			&t.FromLocationID,
			&t.FromLocation,
			&t.ToLocationID,
			&t.ToLocation,
			&t.TransferDate,
			&t.Notes,
		)
		if err != nil {
			return nil, err
		}

		transfers = append(transfers, t)
	}

	return transfers, nil
}

func (r *TransferRepository) GetByLocation(ctx context.Context, locationID int) ([]models.AssetTransfer, error) {
	query := `
		SELECT t.id, t.asset_id, a.name as asset_name,
		       t.employee_id, e.full_name as employee_name,
		       t.from_location_id, fl.address as from_location,
		       t.to_location_id, tl.address as to_location,
		       t.transfer_date, t.notes
		FROM asset_transfers t
		JOIN assets a ON t.asset_id = a.id
		JOIN employees e ON t.employee_id = e.id
		JOIN locations fl ON t.from_location_id = fl.id
		JOIN locations tl ON t.to_location_id = tl.id
		WHERE t.from_location_id = $1 OR t.to_location_id = $1
		ORDER BY t.transfer_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []models.AssetTransfer
	for rows.Next() {
		var t models.AssetTransfer
		err := rows.Scan(
			&t.ID,
			&t.AssetID,
			&t.AssetName,
			&t.EmployeeID,
			&t.EmployeeName,
			&t.FromLocationID,
			&t.FromLocation,
			&t.ToLocationID,
			&t.ToLocation,
			&t.TransferDate,
			&t.Notes,
		)
		if err != nil {
			return nil, err
		}

		transfers = append(transfers, t)
	}

	return transfers, nil
}

func (r *TransferRepository) GetByDateRange(ctx context.Context, from, to time.Time) ([]models.AssetTransfer, error) {
	query := `
		SELECT t.id, t.asset_id, a.name as asset_name,
		       t.employee_id, e.full_name as employee_name,
		       t.from_location_id, fl.address as from_location,
		       t.to_location_id, tl.address as to_location,
		       t.transfer_date, t.notes
		FROM asset_transfers t
		JOIN assets a ON t.asset_id = a.id
		JOIN employees e ON t.employee_id = e.id
		JOIN locations fl ON t.from_location_id = fl.id
		JOIN locations tl ON t.to_location_id = tl.id
		WHERE t.transfer_date BETWEEN $1 AND $2
		ORDER BY t.transfer_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []models.AssetTransfer
	for rows.Next() {
		var t models.AssetTransfer
		err := rows.Scan(
			&t.ID,
			&t.AssetID,
			&t.AssetName,
			&t.EmployeeID,
			&t.EmployeeName,
			&t.FromLocationID,
			&t.FromLocation,
			&t.ToLocationID,
			&t.ToLocation,
			&t.TransferDate,
			&t.Notes,
		)
		if err != nil {
			return nil, err
		}

		transfers = append(transfers, t)
	}

	return transfers, nil
}
