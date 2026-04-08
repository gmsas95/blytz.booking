package services

import (
	"fmt"
	"regexp"
	"strings"

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

	user := models.User{Email: email, Name: name, PasswordHash: hashedPassword}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		business := models.Business{
			Name:        fmt.Sprintf("%s Workshop", name),
			Slug:        buildWorkshopSlug(name),
			Vertical:    "Automotive",
			Description: "Your workshop profile",
			ThemeColor:  "blue",
		}
		if err := tx.Create(&business).Error; err != nil {
			return err
		}

		membership := models.Membership{UserID: user.ID, BusinessID: business.ID, Role: models.MembershipRoleOwner}
		if err := tx.Create(&membership).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, "", err
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func buildWorkshopSlug(name string) string {
	normalized := strings.ToLower(strings.TrimSpace(name))
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug := strings.Trim(re.ReplaceAllString(normalized, "-"), "-")
	if slug == "" {
		slug = "workshop"
	}
	return fmt.Sprintf("%s-%s", slug, strings.ToLower(uuid.NewString()[:6]))
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

func (s *AuthService) GetMemberships(userID uuid.UUID) ([]models.Membership, error) {
	var memberships []models.Membership
	if err := s.DB.Where("user_id = ?", userID).Preload("Business").Order("created_at ASC").Find(&memberships).Error; err != nil {
		return nil, err
	}
	return memberships, nil
}

func (s *AuthService) UserHasBusinessAccess(userID, businessID uuid.UUID) (bool, error) {
	var count int64
	if err := s.DB.Model(&models.Membership{}).Where("user_id = ? AND business_id = ?", userID, businessID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
