package services

import (
	"context"
	"errors"
	"time"

	"blytz.cloud/backend/internal/models"
	"blytz.cloud/backend/internal/repository"
	"blytz.cloud/backend/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtManager       *utils.JWTManager
	db               *gorm.DB
}

func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtManager *utils.JWTManager,
	db *gorm.DB,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
		db:               db,
	}
}

type RegisterRequest struct {
	BusinessName     string `json:"business_name" validate:"required,min=3"`
	BusinessSlug     string `json:"business_slug" validate:"required,alphanum"`
	BusinessVertical string `json:"business_vertical" validate:"required,oneof=automotive wellness creative professional other"`
	Email            string `json:"email" validate:"required,email"`
	Password         string `json:"password" validate:"required,min=8"`
	FullName         string `json:"full_name" validate:"required,min=2"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User         *models.User     `json:"user"`
	Business     *models.Business `json:"business,omitempty"`
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	ExpiresIn    int              `json:"expires_in"` // seconds
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	businessID := uuid.New()

	var authResponse *AuthResponse
	err = s.db.Transaction(func(tx *gorm.DB) error {
		business := &models.Business{
			ID:         businessID,
			Name:       req.BusinessName,
			Slug:       req.BusinessSlug,
			Vertical:   req.BusinessVertical,
			ThemeColor: "blue",
		}
		if err := tx.Create(business).Error; err != nil {
			return err
		}

		user := &models.User{
			ID:            uuid.New(),
			Email:         req.Email,
			PasswordHash:  hashedPassword,
			Name:          req.FullName,
			Role:          models.UserRoleOwner,
			BusinessID:    &businessID,
			IsActive:      true,
			EmailVerified: false,
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, businessID, string(user.Role))
		if err != nil {
			return err
		}

		refreshToken := s.jwtManager.GenerateRefreshToken()
		refreshTokenRecord := &models.RefreshToken{
			ID:        uuid.New(),
			UserID:    user.ID,
			Token:     refreshToken,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		}
		if err := tx.Create(refreshTokenRecord).Error; err != nil {
			return err
		}

		user.PasswordHash = ""
		authResponse = &AuthResponse{
			User:         user,
			Business:     business,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    900, // 15 minutes
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return authResponse, nil
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, errors.New("invalid credentials")
	}

	var business *models.Business
	if user.BusinessID != nil {
		err = s.db.WithContext(ctx).Where("id = ?", *user.BusinessID).First(business).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	var businessID uuid.UUID
	if user.BusinessID != nil {
		businessID = *user.BusinessID
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, businessID, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken := s.jwtManager.GenerateRefreshToken()
	refreshTokenRecord := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshTokenRecord); err != nil {
		return nil, err
	}

	user.PasswordHash = ""

	return &AuthResponse{
		User:         user,
		Business:     business,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	tokenRecord, err := s.refreshTokenRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid refresh token")
		}
		return nil, err
	}

	if tokenRecord.IsExpired() {
		s.refreshTokenRepo.Delete(ctx, refreshToken)
		return nil, errors.New("refresh token expired")
	}

	user, err := s.userRepo.GetByID(ctx, tokenRecord.UserID.String())
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	var business *models.Business
	var businessID uuid.UUID
	if user.BusinessID != nil {
		businessID = *user.BusinessID
		err = s.db.WithContext(ctx).Where("id = ?", businessID).First(business).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, businessID, string(user.Role))
	if err != nil {
		return nil, err
	}

	newRefreshToken := s.jwtManager.GenerateRefreshToken()

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.refreshTokenRepo.WithTx(tx).Delete(ctx, refreshToken); err != nil {
			return err
		}

		newTokenRecord := &models.RefreshToken{
			ID:        uuid.New(),
			UserID:    user.ID,
			Token:     newRefreshToken,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		}
		return s.refreshTokenRepo.WithTx(tx).Create(ctx, newTokenRecord)
	}); err != nil {
		return nil, err
	}

	user.PasswordHash = ""

	return &AuthResponse{
		User:         user,
		Business:     business,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    900,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID) error {
	if err := s.refreshTokenRepo.DeleteByUserID(ctx, userID.String()); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) GetMe(ctx context.Context, userID uuid.UUID) (*models.User, *models.Business, error) {
	user, err := s.userRepo.GetByID(ctx, userID.String())
	if err != nil {
		return nil, nil, err
	}

	var business *models.Business
	if user.BusinessID != nil {
		business = &models.Business{}
		err = s.db.WithContext(ctx).Where("id = ?", *user.BusinessID).First(business).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil, err
			}
			business = nil
		}
	}

	user.PasswordHash = ""
	return user, business, nil
}
