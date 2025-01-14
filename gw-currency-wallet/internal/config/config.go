package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	JWTSecret  string
}

func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load("config.env"); err != nil {
		fmt.Println("Warning: No config.env file found. Using default.")
	}

	return &Config{
		DBUser:     getEnvWithDefault("DB_USER", "default_user"),
		DBPassword: getEnvWithDefault("DB_PASSWORD", "default_password"),
		DBName:     getEnvWithDefault("DB_NAME", "default_db"),
		DBHost:     getEnvWithDefault("DB_HOST", "localhost"),
		DBPort:     getEnvWithDefault("DB_PORT", "5432"),
		JWTSecret:  getEnvWithDefault("JWT_SECRET", "default_jwt_secret"),
	}, nil
}
