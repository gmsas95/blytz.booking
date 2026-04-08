package repository

import (
	"fmt"
	"time"

	"blytz.cloud/backend/config"
	"blytz.cloud/backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
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

	return &Repository{
		DB: db,
	}, nil
}

func (r *Repository) AutoMigrate() error {
	return r.DB.AutoMigrate(
		&models.Business{},
		&models.Service{},
		&models.Slot{},
		&models.Booking{},
		&models.User{},
		&models.Membership{},
		&models.Customer{},
		&models.Vehicle{},
		&models.Job{},
	)
}

func (r *Repository) BackfillMoneyToMinorUnits(defaultCurrencyCode string) error {
	if defaultCurrencyCode == "" {
		defaultCurrencyCode = "USD"
	}

	if r.DB.Migrator().HasColumn(&models.Service{}, "total_price") {
		if err := r.DB.Exec(`
			UPDATE services
			SET total_price_minor = ROUND(total_price * 100),
				deposit_amount_minor = ROUND(deposit_amount * 100),
				currency_code = COALESCE(NULLIF(currency_code, ''), ?)
			WHERE total_price > 0
			  AND (
				total_price_minor = 0
				OR deposit_amount_minor = 0
				OR currency_code IS NULL
				OR currency_code = ''
			  )
		`, defaultCurrencyCode).Error; err != nil {
			return fmt.Errorf("backfill service money fields: %w", err)
		}
	}

	if r.DB.Migrator().HasColumn(&models.Booking{}, "total_price") {
		if err := r.DB.Exec(`
			UPDATE bookings
			SET total_price_minor = ROUND(total_price * 100),
				deposit_paid_minor = ROUND(deposit_paid * 100),
				currency_code = COALESCE(NULLIF(currency_code, ''), ?)
			WHERE total_price > 0
			  AND (
				total_price_minor = 0
				OR deposit_paid_minor = 0
				OR currency_code IS NULL
				OR currency_code = ''
			  )
		`, defaultCurrencyCode).Error; err != nil {
			return fmt.Errorf("backfill booking money fields: %w", err)
		}
	}

	return nil
}

func (r *Repository) SeedData() error {
	// Check if data already exists
	var count int64
	r.DB.Model(&models.Business{}).Count(&count)
	if count > 0 {
		return nil
	}

	// Seed workshops
	businesses := []models.Business{
		{
			Name:        "DetailPro Automotive",
			Slug:        "detail-pro",
			Vertical:    "Automotive",
			Description: "Premium mobile detailing and ceramic coating.",
			ThemeColor:  "blue",
		},
		{
			Name:        "TintLab Studio",
			Slug:        "tint-lab",
			Vertical:    "Automotive",
			Description: "Window tinting and heat rejection packages for daily drivers.",
			ThemeColor:  "emerald",
		},
		{
			Name:        "ShineBay Detailing",
			Slug:        "shine-bay",
			Vertical:    "Automotive",
			Description: "Wash, polish, and protection packages for busy workshop teams.",
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
			BusinessID:         bizList[0].ID,
			Name:               "Full Interior Detail",
			Description:        "Deep clean, steam, shampoo, and leather conditioning.",
			DurationMin:        120,
			TotalPriceMinor:    20000,
			DepositAmountMinor: 5000,
			CurrencyCode:       "USD",
		},
		{
			BusinessID:         bizList[0].ID,
			Name:               "Ceramic Coating Gold",
			Description:        "5-year protection package with paint correction.",
			DurationMin:        360,
			TotalPriceMinor:    120000,
			DepositAmountMinor: 30000,
			CurrencyCode:       "USD",
		},
		{
			BusinessID:         bizList[1].ID,
			Name:               "Nano Ceramic Tint",
			Description:        "Premium tint package with high heat rejection film.",
			DurationMin:        60,
			TotalPriceMinor:    12000,
			DepositAmountMinor: 4000,
			CurrencyCode:       "USD",
		},
		{
			BusinessID:         bizList[1].ID,
			Name:               "Front Two Windows Tint",
			Description:        "Quick tint installation for front window pairs.",
			DurationMin:        45,
			TotalPriceMinor:    18000,
			DepositAmountMinor: 6000,
			CurrencyCode:       "USD",
		},
	}

	for _, svc := range services {
		r.DB.Create(&svc)
	}

	customers := []models.Customer{
		{
			BusinessID: bizList[0].ID,
			Name:       "Alice Smith",
			Email:      "alice@example.com",
			Phone:      "555-0101",
			Notes:      "Prefers morning drop-off.",
		},
		{
			BusinessID: bizList[0].ID,
			Name:       "Marco Rivera",
			Email:      "marco@example.com",
			Phone:      "555-0120",
			Notes:      "Repeat ceramic coating customer.",
		},
	}

	for _, customer := range customers {
		r.DB.Create(&customer)
	}

	var customerList []models.Customer
	r.DB.Where("business_id = ?", bizList[0].ID).Find(&customerList)
	if len(customerList) >= 2 {
		vehicles := []models.Vehicle{
			{
				BusinessID:   bizList[0].ID,
				CustomerID:   customerList[0].ID,
				Year:         2022,
				Make:         "Tesla",
				Model:        "Model Y",
				Color:        "White",
				LicensePlate: "BLYTZ01",
			},
			{
				BusinessID:   bizList[0].ID,
				CustomerID:   customerList[1].ID,
				Year:         2021,
				Make:         "BMW",
				Model:        "X5",
				Color:        "Black",
				LicensePlate: "SHINE22",
			},
		}

		for _, vehicle := range vehicles {
			r.DB.Create(&vehicle)
		}

		var vehicleList []models.Vehicle
		r.DB.Where("business_id = ?", bizList[0].ID).Find(&vehicleList)
		if len(vehicleList) >= 2 {
			jobs := []models.Job{
				{
					BusinessID:  bizList[0].ID,
					CustomerID:  customerList[0].ID,
					VehicleID:   vehicleList[0].ID,
					Title:       "Interior detail and hand wash",
					Status:      models.JobStatusScheduled,
					ScheduledAt: time.Now().UTC().Add(24 * time.Hour),
					Notes:       "Customer waiting in lounge.",
				},
				{
					BusinessID:  bizList[0].ID,
					CustomerID:  customerList[1].ID,
					VehicleID:   vehicleList[1].ID,
					Title:       "Ceramic coating maintenance",
					Status:      models.JobStatusInProgress,
					ScheduledAt: time.Now().UTC().Add(2 * time.Hour),
					Notes:       "Second layer inspection before delivery.",
				},
			}

			for _, job := range jobs {
				r.DB.Create(&job)
			}
		}
	}

	return nil
}
