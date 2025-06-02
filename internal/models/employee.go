package models

type Employee struct {
	ID           int     `json:"id"`
	FullName     string  `json:"full_name"`
	Position     string  `json:"position"`
	Email        string  `json:"email"`
	PasswordHash string  `json:"-"`
	Role         string  `json:"role"`
	DepartmentID *int    `json:"department_id,omitempty"`
	Department   *string `json:"department,omitempty"`
}
