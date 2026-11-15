package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName         string
	AppEnv          string
	AppPort         string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
	JWTSecret       string
	JWTExpiresHours int
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	expHours, err := strconv.Atoi(getEnv("JWT_EXPIRES_HOURS", "24"))
	if err != nil {
		expHours = 24
	}

	return &Config{
		AppName:         getEnv("APP_NAME", "findit-campus-backend"),
		AppEnv:          getEnv("APP_ENV", "development"),
		AppPort:         getEnv("APP_PORT", "8080"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "findit_campus"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
		JWTSecret:       getEnv("JWT_SECRET", "default_dev_secret_change_me"),
		JWTExpiresHours: expHours,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}