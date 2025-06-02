package repository

import (
	"context"
	"database/sql"
	"inventory-system/internal/models"
	"time"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.Employee, error) {
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
		if err == sql.ErrNoRows {
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

func (r *AuthRepository) CreateSession(ctx context.Context, employeeID int, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO sessions (employee_id, token, expires_at)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		employeeID,
		token,
		expiresAt,
	)

	return err
}

func (r *AuthRepository) GetSession(ctx context.Context, token string) (*models.Employee, error) {
	query := `
		SELECT e.id, e.full_name, e.position, e.email, e.role, e.department_id
		FROM sessions s
		JOIN employees e ON s.employee_id = e.id
		WHERE s.token = $1 AND s.expires_at > NOW()
	`

	var e models.Employee
	var deptID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&e.ID,
		&e.FullName,
		&e.Position,
		&e.Email,
		&e.Role,
		&deptID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
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

func (r *AuthRepository) DeleteSession(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM sessions WHERE token = $1",
		token,
	)
	return err
}
