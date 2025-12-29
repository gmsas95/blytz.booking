package repository

import (
	"context"
	"fmt"

	"blytz.cloud/backend/config"
	"blytz.cloud/backend/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	// Connect to PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Repository{
		DB:    db,
		Redis: rdb,
	}, nil
}

func (r *Repository) AutoMigrate() error {
	return r.DB.AutoMigrate(
		&models.Business{},
		&models.Service{},
		&models.Slot{},
		&models.Booking{},
		&models.User{},
		&models.Customer{},
		&models.Payment{},
		&models.Subscription{},
		&models.BusinessSettings{},
		&models.BookingHistory{},
		&models.RecurringSchedule{},
		&models.Webhook{},
		&models.RefreshToken{},
	)
}

func (r *Repository) SeedData() error {
	// Check if data already exists
	var count int64
	r.DB.Model(&models.Business{}).Count(&count)
	if count > 0 {
		return nil
	}

	// Seed businesses
	businesses := []models.Business{
		{
			Name:        "DetailPro Automotive",
			Slug:        "detail-pro",
			Vertical:    "Automotive",
			Description: "Premium mobile detailing and ceramic coating.",
			ThemeColor:  "blue",
		},
		{
			Name:        "Lumina Wellness Spa",
			Slug:        "lumina-spa",
			Vertical:    "Wellness",
			Description: "Massage therapy, facials, and holistic healing.",
			ThemeColor:  "emerald",
		},
		{
			Name:        "FlashFrame Studios",
			Slug:        "flash-frame",
			Vertical:    "Creative",
			Description: "Editorial portraiture and high-end fashion photography.",
			ThemeColor:  "zinc",
		},
	}

	for _, biz := range businesses {
		r.DB.Create(&biz)
	}

	// Get businesses for services/slots reference
	var bizList []models.Business
	r.DB.Find(&bizList)

	if len(bizList) == 0 {
		return nil
	}

	// Seed services
	services := []models.Service{
		{
			BusinessID:    bizList[0].ID,
			Name:          "Full Interior Detail",
			Description:   "Deep clean, steam, shampoo, and leather conditioning.",
			DurationMin:   120,
			TotalPrice:    200,
			DepositAmount: 50,
		},
		{
			BusinessID:    bizList[0].ID,
			Name:          "Ceramic Coating Gold",
			Description:   "5-year protection package with paint correction.",
			DurationMin:   360,
			TotalPrice:    1200,
			DepositAmount: 300,
		},
		{
			BusinessID:    bizList[1].ID,
			Name:          "Deep Tissue Massage",
			Description:   "60-minute therapeutic massage for stress relief.",
			DurationMin:   60,
			TotalPrice:    120,
			DepositAmount: 40,
		},
		{
			BusinessID:    bizList[1].ID,
			Name:          "Hydrafacial Signature",
			Description:   "Cleanse, extract, and hydrate skin.",
			DurationMin:   45,
			TotalPrice:    180,
			DepositAmount: 60,
		},
	}

	for _, svc := range services {
		r.DB.Create(&svc)
	}

	return nil
}
