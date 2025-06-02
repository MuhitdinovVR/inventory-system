package services

import (
	"context"
	"errors"
	"inventory-system/internal/models"
	"inventory-system/internal/repository"
)

var (
	ErrDepartmentNotFound     = errors.New("department not found")
	ErrDepartmentHasEmployees = errors.New("cannot delete department with employees")
)

type DepartmentService struct {
	departmentRepo *repository.DepartmentRepository
	employeeRepo   *repository.EmployeeRepository
}

func NewDepartmentService(
	departmentRepo *repository.DepartmentRepository,
	employeeRepo *repository.EmployeeRepository,
) *DepartmentService {
	return &DepartmentService{
		departmentRepo: departmentRepo,
		employeeRepo:   employeeRepo,
	}
}

func (s *DepartmentService) GetAllDepartments(ctx context.Context) ([]models.Department, error) {
	return s.departmentRepo.GetAll(ctx)
}

func (s *DepartmentService) GetDepartmentByID(ctx context.Context, id int) (*models.Department, error) {
	department, err := s.departmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if department == nil {
		return nil, ErrDepartmentNotFound
	}
	return department, nil
}

func (s *DepartmentService) CreateDepartment(ctx context.Context, department models.Department) (int, error) {
	if department.HeadID != nil {
		exists, err := s.employeeRepo.Exists(ctx, *department.HeadID)
		if err != nil {
			return 0, err
		}
		if !exists {
			return 0, ErrEmployeeNotFound
		}
	}

	return s.departmentRepo.Create(ctx, department)
}

func (s *DepartmentService) UpdateDepartment(ctx context.Context, department models.Department) error {
	if department.HeadID != nil {
		exists, err := s.employeeRepo.Exists(ctx, *department.HeadID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrEmployeeNotFound
		}
	}

	return s.departmentRepo.Update(ctx, department)
}

func (s *DepartmentService) DeleteDepartment(ctx context.Context, id int) error {
	employees, err := s.departmentRepo.GetEmployees(ctx, id)
	if err != nil {
		return err
	}

	if len(employees) > 0 {
		return ErrDepartmentHasEmployees
	}

	return s.departmentRepo.Delete(ctx, id)
}

func (s *DepartmentService) GetEmployeesByDepartment(ctx context.Context, departmentID int) ([]models.Employee, error) {
	return s.departmentRepo.GetEmployees(ctx, departmentID)
}

func (s *DepartmentService) GetPotentialHeads(ctx context.Context) ([]models.Employee, error) {
	return s.employeeRepo.GetAll(ctx)
}

func (s *DepartmentService) IsDepartmentHead(ctx context.Context, employeeID int) (bool, bool) {
	var isHead bool
	err := s.departmentRepo.IsEmployeeDepartmentHead(ctx, employeeID, &isHead)
	return isHead, err
}
