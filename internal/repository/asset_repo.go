package repository

import (
	"context"
	"database/sql"
	"errors"
	"inventory-system/internal/models"
)

type AssetRepository struct {
	db *sql.DB
}

func NewAssetRepository(db *sql.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) GetAll(ctx context.Context) ([]models.Asset, error) {
	query := `
		SELECT a.id, a.name, a.category, a.acquisition_date, a.cost,
		       a.status_id, s.name as status_name,
		       a.current_location_id, l.address as location_address,
		       a.department_id, d.name as department_name
		FROM assets a
		JOIN asset_statuses s ON a.status_id = s.id
		JOIN locations l ON a.current_location_id = l.id
		LEFT JOIN departments d ON a.department_id = d.id
		ORDER BY a.name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var a models.Asset
		var deptID sql.NullInt64
		var deptName sql.NullString

		err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Category,
			&a.AcquisitionDate,
			&a.Cost,
			&a.StatusID,
			&a.Status,
			&a.CurrentLocationID,
			&a.Location,
			&deptID,
			&deptName,
		)
		if err != nil {
			return nil, err
		}

		if deptID.Valid {
			id := int(deptID.Int64)
			a.DepartmentID = &id
			name := deptName.String
			a.Department = &name
		}

		assets = append(assets, a)
	}

	return assets, nil
}

func (r *AssetRepository) GetByID(ctx context.Context, id int) (*models.Asset, error) {
	query := `
		SELECT a.id, a.name, a.category, a.acquisition_date, a.cost,
		       a.status_id, s.name as status_name,
		       a.current_location_id, l.address as location_address,
		       a.department_id, d.name as department_name
		FROM assets a
		JOIN asset_statuses s ON a.status_id = s.id
		JOIN locations l ON a.current_location_id = l.id
		LEFT JOIN departments d ON a.department_id = d.id
		WHERE a.id = $1
	`

	var a models.Asset
	var deptID sql.NullInt64
	var deptName sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID,
		&a.Name,
		&a.Category,
		&a.AcquisitionDate,
		&a.Cost,
		&a.StatusID,
		&a.Status,
		&a.CurrentLocationID,
		&a.Location,
		&deptID,
		&deptName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if deptID.Valid {
		id := int(deptID.Int64)
		a.DepartmentID = &id
		name := deptName.String
		a.Department = &name
	}

	return &a, nil
}

func (r *AssetRepository) Create(ctx context.Context, asset models.Asset) (int, error) {
	query := `
		INSERT INTO assets (name, category, acquisition_date, cost, 
		                  status_id, current_location_id, department_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(
		ctx,
		query,
		asset.Name,
		asset.Category,
		asset.AcquisitionDate,
		asset.Cost,
		asset.StatusID,
		asset.CurrentLocationID,
		asset.DepartmentID,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AssetRepository) Update(ctx context.Context, asset models.Asset) error {
	query := `
		UPDATE assets
		SET name = $1, category = $2, acquisition_date = $3, cost = $4,
		    status_id = $5, current_location_id = $6, department_id = $7
		WHERE id = $8
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		asset.Name,
		asset.Category,
		asset.AcquisitionDate,
		asset.Cost,
		asset.StatusID,
		asset.CurrentLocationID,
		asset.DepartmentID,
		asset.ID,
	)

	return err
}

func (r *AssetRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM assets WHERE id = $1", id)
	return err
}

func (r *AssetRepository) GetByStatus(ctx context.Context, statusID int) ([]models.Asset, error) {
	query := `
		SELECT a.id, a.name, a.category, a.acquisition_date, a.cost,
		       a.status_id, s.name as status_name,
		       a.current_location_id, l.address as location_address,
		       a.department_id, d.name as department_name
		FROM assets a
		JOIN asset_statuses s ON a.status_id = s.id
		JOIN locations l ON a.current_location_id = l.id
		LEFT JOIN departments d ON a.department_id = d.id
		WHERE a.status_id = $1
		ORDER BY a.name
	`

	rows, err := r.db.QueryContext(ctx, query, statusID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var a models.Asset
		var deptID sql.NullInt64
		var deptName sql.NullString

		err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Category,
			&a.AcquisitionDate,
			&a.Cost,
			&a.StatusID,
			&a.Status,
			&a.CurrentLocationID,
			&a.Location,
			&deptID,
			&deptName,
		)
		if err != nil {
			return nil, err
		}

		if deptID.Valid {
			id := int(deptID.Int64)
			a.DepartmentID = &id
			name := deptName.String
			a.Department = &name
		}

		assets = append(assets, a)
	}

	return assets, nil
}

func (r *AssetRepository) GetByLocation(ctx context.Context, locationID int) ([]models.Asset, error) {
	query := `
		SELECT a.id, a.name, a.category, a.acquisition_date, a.cost,
		       a.status_id, s.name as status_name,
		       a.current_location_id, l.address as location_address,
		       a.department_id, d.name as department_name
		FROM assets a
		JOIN asset_statuses s ON a.status_id = s.id
		JOIN locations l ON a.current_location_id = l.id
		LEFT JOIN departments d ON a.department_id = d.id
		WHERE a.current_location_id = $1
		ORDER BY a.name
	`

	rows, err := r.db.QueryContext(ctx, query, locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var a models.Asset
		var deptID sql.NullInt64
		var deptName sql.NullString

		err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Category,
			&a.AcquisitionDate,
			&a.Cost,
			&a.StatusID,
			&a.Status,
			&a.CurrentLocationID,
			&a.Location,
			&deptID,
			&deptName,
		)
		if err != nil {
			return nil, err
		}

		if deptID.Valid {
			id := int(deptID.Int64)
			a.DepartmentID = &id
			name := deptName.String
			a.Department = &name
		}

		assets = append(assets, a)
	}

	return assets, nil
}

func (r *AssetRepository) GetByDepartment(ctx context.Context, departmentID int) ([]models.Asset, error) {
	query := `
		SELECT a.id, a.name, a.category, a.acquisition_date, a.cost,
		       a.status_id, s.name as status_name,
		       a.current_location_id, l.address as location_address,
		       a.department_id, d.name as department_name
		FROM assets a
		JOIN asset_statuses s ON a.status_id = s.id
		JOIN locations l ON a.current_location_id = l.id
		LEFT JOIN departments d ON a.department_id = d.id
		WHERE a.department_id = $1
		ORDER BY a.name
	`

	rows, err := r.db.QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var a models.Asset
		var deptID sql.NullInt64
		var deptName sql.NullString

		err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Category,
			&a.AcquisitionDate,
			&a.Cost,
			&a.StatusID,
			&a.Status,
			&a.CurrentLocationID,
			&a.Location,
			&deptID,
			&deptName,
		)
		if err != nil {
			return nil, err
		}

		if deptID.Valid {
			id := int(deptID.Int64)
			a.DepartmentID = &id
			name := deptName.String
			a.Department = &name
		}

		assets = append(assets, a)
	}

	return assets, nil
}

func (r *AssetRepository) UpdateStatus(ctx context.Context, id int, statusID int) error {
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE assets SET status_id = $1 WHERE id = $2",
		statusID,
		id,
	)
	return err
}

func (r *AssetRepository) UpdateLocation(ctx context.Context, id int, locationID int) error {
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE assets SET current_location_id = $1 WHERE id = $2",
		locationID,
		id,
	)
	return err
}

func (r *AssetRepository) GetTransferHistory(ctx context.Context, assetID int) ([]models.AssetTransfer, error) {
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

func (r *AssetRepository) Exists(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM assets WHERE id = $1)",
		id,
	).Scan(&exists)

	return exists, err
}

func (r *AssetRepository) HasTransfers(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM asset_transfers WHERE asset_id = $1)",
		id,
	).Scan(&exists)

	return exists, err
}
