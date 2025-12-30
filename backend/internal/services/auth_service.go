package services

import (
	"blytz.cloud/backend/internal/auth"
	"blytz.cloud/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	*BaseService
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		BaseService: NewBaseService(db),
	}
}

func (s *AuthService) Register(email, name, password string) (*models.User, string, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, "", ErrConflict
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	// Create user
	user := models.User{
		Email:        email,
		Name:         name,
		PasswordHash: hashedPassword,
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, "", err
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	// Find user
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", ErrUnauthorized
		}
		return nil, "", err
	}

	// Check password
	if !auth.CheckPassword(password, user.PasswordHash) {
		return nil, "", ErrUnauthorized
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *AuthService) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}
