package services

import (
	"crypto/rand"
	"encoding/hex"
	"time"

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

func (s *AuthService) ForgotPassword(email string) (string, error) {
	// Find user
	user, err := s.GetByEmail(email)
	if err != nil {
		return "", ErrNotFound
	}

	// Generate reset token
	token := generateRandomToken()

	// Save token and expiration (1 hour)
	now := time.Now()
	expires := now.Add(1 * time.Hour)

	if err := s.DB.Model(user).Updates(map[string]interface{}{
		"password_reset_token": token,
		"password_reset_at":    expires,
	}).Error; err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	// Find user with valid token
	var user models.User
	err := s.DB.Where("password_reset_token = ? AND password_reset_at > ?", token, time.Now()).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrBadRequest
		}
		return err
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password and clear token
	return s.DB.Model(&user).Updates(map[string]interface{}{
		"password_hash":        hashedPassword,
		"password_reset_token": nil,
		"password_reset_at":    nil,
	}).Error
}

func generateRandomToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
