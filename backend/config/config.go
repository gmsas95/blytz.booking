package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Email    EmailConfig
}

type ServerConfig struct {
	Port       string
	Env        string
	BaseDomain string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret string
}

type EmailConfig struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:       getEnv("SERVER_PORT", "8080"),
			Env:        getEnv("ENV", "development"),
			BaseDomain: getEnv("BASE_DOMAIN", "blytz.cloud"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "blytz"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "change-me-in-production"),
		},
		Email: EmailConfig{
			From:     getEnv("EMAIL_FROM", "noreply@blytz.cloud"),
			Host:     getEnv("EMAIL_HOST", "smtp.gmail.com"),
			Port:     getEnvAsInt("EMAIL_PORT", 587),
			Username: getEnv("EMAIL_USERNAME", ""),
			Password: getEnv("EMAIL_PASSWORD", ""),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}
