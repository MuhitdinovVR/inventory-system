package repository

import (
	"context"
	"database/sql"
	"errors"
	"inventory-system/internal/models"
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) GetAll(ctx context.Context) ([]models.Employee, error) {
	query := `
		SELECT e.id, e.full_name, e.position, e.email, e.role, 
		       e.department_id, d.name as department_name
		FROM employees e
		LEFT JOIN departments d ON e.department_id = d.id
		ORDER BY e.full_name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var e models.Employee
		var deptID sql.NullInt64
		var deptName sql.NullString

		err := rows.Scan(
			&e.ID,
			&e.FullName,
			&e.Position,
			&e.Email,
			&e.Role,
			&deptID,
			&deptName,
		)
		if err != nil {
			return nil, err
		}

		if deptID.Valid {
			id := int(deptID.Int64)
			e.DepartmentID = &id
			name := deptName.String
			e.Department = &name
		}

		employees = append(employees, e)
	}

	return employees, nil
}

func (r *EmployeeRepository) GetByID(ctx context.Context, id int) (*models.Employee, error) {
	query := `
		SELECT e.id, e.full_name, e.position, e.email, e.role, 
		       e.department_id, d.name as department_name
		FROM employees e
		LEFT JOIN departments d ON e.department_id = d.id
		WHERE e.id = $1
	`

	var e models.Employee
	var deptID sql.NullInt64
	var deptName sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID,
		&e.FullName,
		&e.Position,
		&e.Email,
		&e.Role,
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
		e.DepartmentID = &id
		name := deptName.String
		e.Department = &name
	}

	return &e, nil
}

func (r *EmployeeRepository) GetByEmail(ctx context.Context, email string) (*models.Employee, error) {
	query := `
		SELECT id, full_name, position, email, password_hash, role, department_id
		FROM employees
		WHERE email = $1
	`

	var e models.Employee
	var deptID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&e.ID,
		&e.FullName,
		&e.Position,
		&e.Email,
		&e.PasswordHash,
		&e.Role,
		&deptID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if deptID.Valid {
		id := int(deptID.Int64)
		e.DepartmentID = &id
	}

	return &e, nil
}

func (r *EmployeeRepository) Create(ctx context.Context, employee models.Employee) (int, error) {
	query := `
		INSERT INTO employees (full_name, position, email, password_hash, role, department_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(
		ctx,
		query,
		employee.FullName,
		employee.Position,
		employee.Email,
		employee.PasswordHash,
		employee.Role,
		employee.DepartmentID,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *EmployeeRepository) Update(ctx context.Context, employee models.Employee) error {
	query := `
		UPDATE employees
		SET full_name = $1, position = $2, email = $3, 
		    role = $4, department_id = $5
		WHERE id = $6
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		employee.FullName,
		employee.Position,
		employee.Email,
		employee.Role,
		employee.DepartmentID,
		employee.ID,
	)

	return err
}

func (r *EmployeeRepository) UpdatePassword(ctx context.Context, id int, passwordHash string) error {
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE employees SET password_hash = $1 WHERE id = $2",
		passwordHash,
		id,
	)
	return err
}

func (r *EmployeeRepository) Delete(ctx context.Context, id int) error {
	// Проверяем, является ли сотрудник главой отдела
	var isHead bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM departments WHERE head_id = $1)",
		id,
	).Scan(&isHead)

	if err != nil {
		return err
	}

	if isHead {
		return errors.New("cannot delete employee who is a department head")
	}

	_, err = r.db.ExecContext(ctx, "DELETE FROM employees WHERE id = $1", id)
	return err
}

func (r *EmployeeRepository) Exists(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM employees WHERE id = $1)",
		id,
	).Scan(&exists)

	return exists, err
}

func (r *EmployeeRepository) IsEmployeeDepartmentHead(ctx context.Context, id int) (bool, error) {
	var isHead bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM departments WHERE head_id = $1)",
		id,
	).Scan(&isHead)

	return isHead, err
}
