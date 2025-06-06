package services

import (
	"context"
	"errors"
	"inventory-system/internal/models"
	"inventory-system/internal/repository"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	employeeRepo *repository.EmployeeRepository
	secretKey    []byte
	tokenExpiry  time.Duration
}

func NewAuthService(
	employeeRepo *repository.EmployeeRepository,
	secretKey string,
	tokenExpiry time.Duration,
) *AuthService {
	return &AuthService{
		employeeRepo: employeeRepo,
		secretKey:    []byte(secretKey),
		tokenExpiry:  tokenExpiry,
	}
}
func (s *AuthService) Authenticate(ctx context.Context, email, password string) (*models.Employee, error) {
	log.Printf("Authenticating user: %s", email)
	employee, err := s.employeeRepo.GetByEmail(ctx, email)
	if err != nil {
		log.Printf("Database error: %v", err)
		return nil, err
	}
	if employee == nil {
		log.Printf("User not found: %s", email)
		return nil, ErrInvalidCredentials
	}

	log.Printf("Stored hash: %s", employee.PasswordHash)

	// Сравнение пароля с хешем
	err = bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte(password))
	if err != nil {
		log.Printf("Password mismatch: %v", err)

		// Дополнительная отладочная информация
		log.Printf("Input password: %s", password)
		log.Printf("Hash length: %d", len(employee.PasswordHash))

		return nil, ErrInvalidCredentials
	}

	return employee, nil
}

func (s *AuthService) GenerateToken(employee *models.Employee) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.tokenExpiry)

	claims := jwt.MapClaims{
		"sub":  employee.ID,
		"exp":  expiresAt.Unix(),
		"role": employee.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*models.Employee, error) {
	// Удаляем префикс "Bearer" если он есть
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		employeeID := int(claims["sub"].(float64))
		employee, err := s.employeeRepo.GetByID(ctx, employeeID)
		if err != nil {
			return nil, err
		}
		if employee == nil {
			return nil, ErrInvalidToken
		}
		return employee, nil
	}

	return nil, ErrInvalidToken
}

func (s *AuthService) InvalidateToken(tokenString string) error {
	// В реальном приложении здесь можно добавить токен в черный список
	return nil
}
