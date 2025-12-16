package config

import (
	"os"
)

type Config struct {
	AppName  string
	AppEnv   string
	AppDebug bool
	AppHost  string
	AppPort  string

	DBHost     string
	DBPort     string
	DBDatabase string
	DBUsername string
	DBPassword string

	JWTSecret     string
	JWTTTL        int
	JWTRefreshTTL int
}

func Load() *Config {
	return &Config{
		AppName:  getEnv("APP_NAME", "DokumenKeuangan"),
		AppEnv:   getEnv("APP_ENV", "local"),
		AppDebug: getEnv("APP_DEBUG", "true") == "true",
		AppHost:  getEnv("APP_HOST", "0.0.0.0"),
		AppPort:  getEnv("APP_PORT", "8000"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBDatabase: getEnv("DB_DATABASE", "dokumen_keuangan"),
		DBUsername: getEnv("DB_USERNAME", "dokumen_user"),
		DBPassword: getEnv("DB_PASSWORD", "dokumen_pass"),

		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		JWTTTL:        60,
		JWTRefreshTTL: 20160,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
