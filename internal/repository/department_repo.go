package repository

import (
	"context"
	"database/sql"
	"errors"
	"inventory-system/internal/models"
)

type DepartmentRepository struct {
	db *sql.DB
}

func NewDepartmentRepository(db *sql.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) GetAll(ctx context.Context) ([]models.Department, error) {
	query := `
		SELECT d.id, d.name, d.location, 
		       d.head_id, e.full_name as head_name
		FROM departments d
		LEFT JOIN employees e ON d.head_id = e.id
		ORDER BY d.name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []models.Department
	for rows.Next() {
		var d models.Department
		var headID sql.NullInt64
		var headName sql.NullString

		err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Location,
			&headID,
			&headName,
		)
		if err != nil {
			return nil, err
		}

		if headID.Valid {
			id := int(headID.Int64)
			d.HeadID = &id
			name := headName.String
			d.HeadName = &name
		}

		departments = append(departments, d)
	}

	return departments, nil
}

func (r *DepartmentRepository) GetByID(ctx context.Context, id int) (*models.Department, error) {
	query := `
		SELECT d.id, d.name, d.location, 
		       d.head_id, e.full_name as head_name
		FROM departments d
		LEFT JOIN employees e ON d.head_id = e.id
		WHERE d.id = $1
	`

	var d models.Department
	var headID sql.NullInt64
	var headName sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID,
		&d.Name,
		&d.Location,
		&headID,
		&headName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if headID.Valid {
		id := int(headID.Int64)
		d.HeadID = &id
		name := headName.String
		d.HeadName = &name
	}

	return &d, nil
}

func (r *DepartmentRepository) Create(ctx context.Context, department models.Department) (int, error) {
	query := `
		INSERT INTO departments (name, location, head_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(
		ctx,
		query,
		department.Name,
		department.Location,
		department.HeadID,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *DepartmentRepository) Update(ctx context.Context, department models.Department) error {
	query := `
		UPDATE departments
		SET name = $1, location = $2, head_id = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		department.Name,
		department.Location,
		department.HeadID,
		department.ID,
	)

	return err
}

func (r *DepartmentRepository) Delete(ctx context.Context, id int) error {
	// Проверяем, есть ли сотрудники в отделе
	var hasEmployees bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM employees WHERE department_id = $1)",
		id,
	).Scan(&hasEmployees)

	if err != nil {
		return err
	}

	if hasEmployees {
		return errors.New("cannot delete department with employees")
	}

	_, err = r.db.ExecContext(ctx, "DELETE FROM departments WHERE id = $1", id)
	return err
}

func (r *DepartmentRepository) GetEmployees(ctx context.Context, departmentID int) ([]models.Employee, error) {
	query := `
		SELECT id, full_name, position, email, role
		FROM employees
		WHERE department_id = $1
		ORDER BY full_name
	`

	rows, err := r.db.QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var e models.Employee
		err := rows.Scan(
			&e.ID,
			&e.FullName,
			&e.Position,
			&e.Email,
			&e.Role,
		)
		if err != nil {
			return nil, err
		}

		employees = append(employees, e)
	}

	return employees, nil
}

func (r *DepartmentRepository) Exists(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM departments WHERE id = $1)",
		id,
	).Scan(&exists)

	return exists, err
}

func (r *DepartmentRepository) IsEmployeeDepartmentHead(ctx context.Context, employeeID int, b *bool) bool {
	var isHead bool
	_ = r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM departments WHERE head_id = $1)",
		employeeID,
	).Scan(&isHead)

	return isHead
}
