package services

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"inventory-system/internal/models"
	"inventory-system/internal/repository"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")

	ErrEmployeeNotFound = errors.New("employee not found")
)

type EmployeeService struct {
	employeeRepo   *repository.EmployeeRepository
	departmentRepo *repository.DepartmentRepository
}

func NewEmployeeService(
	employeeRepo *repository.EmployeeRepository,
	departmentRepo *repository.DepartmentRepository,
) *EmployeeService {
	return &EmployeeService{
		employeeRepo:   employeeRepo,
		departmentRepo: departmentRepo,
	}
}

func (s *EmployeeService) GetAllEmployees(ctx context.Context) ([]models.Employee, error) {
	return s.employeeRepo.GetAll(ctx)
}

func (s *EmployeeService) GetEmployeeByID(ctx context.Context, id int) (*models.Employee, error) {
	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, ErrEmployeeNotFound
	}
	return employee, nil
}

func (s *EmployeeService) CreateEmployee(ctx context.Context, employee models.Employee) (int, error) {
	// Validate department if specified
	if employee.DepartmentID != nil {
		exists, err := s.departmentRepo.Exists(ctx, *employee.DepartmentID)
		if err != nil {
			return 0, err
		}
		if !exists {
			return 0, ErrDepartmentNotFound
		}
	}

	// Check if email already exists
	existing, err := s.employeeRepo.GetByEmail(ctx, employee.Email)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(employee.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	employee.PasswordHash = string(hashedPassword)

	return s.employeeRepo.Create(ctx, employee)
}

func (s *EmployeeService) UpdateEmployee(ctx context.Context, employee models.Employee) error {
	// Validate department if specified
	if employee.DepartmentID != nil {
		exists, err := s.departmentRepo.Exists(ctx, *employee.DepartmentID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrDepartmentNotFound
		}
	}

	// Check if email belongs to another employee
	existing, err := s.employeeRepo.GetByEmail(ctx, employee.Email)
	if err != nil {
		return err
	}
	if existing != nil && existing.ID != employee.ID {
		return ErrEmailAlreadyExists
	}

	return s.employeeRepo.Update(ctx, employee)
}

func (s *EmployeeService) DeleteEmployee(ctx context.Context, id int) error {
	// Check if employee is a department head
	isHead, err := s.employeeRepo.IsEmployeeDepartmentHead(ctx, id)
	if err != nil {
		return err
	}
	if isHead {
		return errors.New("cannot delete employee who is a department head")
	}

	return s.employeeRepo.Delete(ctx, id)
}

func (s *EmployeeService) Authenticate(ctx context.Context, email, password string) (*models.Employee, error) {
	employee, err := s.employeeRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return employee, nil
}

func (s *EmployeeService) UpdatePassword(ctx context.Context, id int, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.employeeRepo.UpdatePassword(ctx, id, string(hashedPassword))
}
