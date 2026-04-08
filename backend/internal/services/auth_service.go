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

const dummyPasswordHash = "$2a$10$7EqJtq98hPqEX7fNZaFWoO.HxQ9gZQh1g0X1p1rRZ8bG8z2u4Vt6G"

type AuthService struct {
	*BaseService
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		BaseService: NewBaseService(db),
	}
}

func (s *AuthService) Register(email, name, password string) (*models.User, string, error) {
	var existingUser models.User
	if err := s.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		auth.CheckPassword(password, dummyPasswordHash)
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
	token, err := auth.GenerateToken(user.ID.String(), user.Email, user.TokenVersion)
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
			auth.CheckPassword(password, dummyPasswordHash)
			return nil, "", ErrUnauthorized
		}
		return nil, "", err
	}

	// Check password
	if !auth.CheckPassword(password, user.PasswordHash) {
		return nil, "", ErrUnauthorized
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID.String(), user.Email, user.TokenVersion)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *AuthService) ValidateUserSession(userID uuid.UUID, tokenVersion int) (*models.User, error) {
	user, err := s.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user.TokenVersion != tokenVersion {
		return nil, ErrUnauthorized
	}
	return user, nil
}

func (s *AuthService) RevokeUserSessions(userID uuid.UUID) error {
	result := s.DB.Model(&models.User{}).Where("id = ?", userID).Update("token_version", gorm.Expr("token_version + 1"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
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
