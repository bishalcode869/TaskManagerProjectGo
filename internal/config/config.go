package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// AppConfig holds all the configuration settings
type AppConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	JWTSecret  string
}

var Config *AppConfig

// LoadConfig loads the environment variables and sets them in Config
func LoadConfig() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// Try loading from parent dir if not found in current dir
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("⚠️  No .env file found. Using system envs.")
		}
	}

	Config = &AppConfig{
		DBHost:     mustGetEnvOrDefault("DB_HOST", "localhost"),
		DBUser:     mustGetEnvOrDefault("DB_USER", "postgres"),
		DBPassword: mustGetEnvOrDefault("DB_PASSWORD", "bishal1212"),
		DBName:     mustGetEnvOrDefault("DB_NAME", "gotasker"),
		DBPort:     mustGetEnvOrDefault("DB_PORT", "5432"),
		JWTSecret:  mustGetEnvOrDefault("JWT_SECRET", "mySuperSecretKey"),
	}

	log.Println("✅ Configuration loaded successfully.")
}

// mustGetEnvOrDefault returns env var or a default value
func mustGetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
