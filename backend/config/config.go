package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	CORS     CORSConfig
	Startup  StartupConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type CORSConfig struct {
	AllowedOrigins []string
}

type StartupConfig struct {
	AutoMigrate   bool
	SeedData      bool
	BackfillMoney bool
}

type JWTConfig struct {
	Secret         string
	CookieName     string
	ForceSecure    bool
	TrustedProxies []string
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "blytz"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173"),
		},
		Startup: StartupConfig{
			AutoMigrate:   getEnvAsBool("AUTO_MIGRATE", getEnv("ENV", "development") != "production"),
			SeedData:      getEnvAsBool("SEED_DATA", false),
			BackfillMoney: getEnvAsBool("BACKFILL_MONEY_FIELDS", getEnv("ENV", "development") != "production"),
		},
		JWT: JWTConfig{
			Secret:         getEnv("JWT_SECRET", ""),
			CookieName:     getEnv("JWT_COOKIE_NAME", "blytz_session"),
			ForceSecure:    getEnvAsBool("JWT_COOKIE_SECURE", getEnv("ENV", "development") == "production"),
			TrustedProxies: getEnvAsSlice("TRUSTED_PROXIES", "127.0.0.1"),
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

func getEnvAsBool(key string, defaultVal bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultVal
}

func getEnvAsSlice(key, defaultVal string) []string {
	value := getEnv(key, defaultVal)
	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	return origins
}
