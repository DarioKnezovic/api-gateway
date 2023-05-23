package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

// Config holds the configuration values
type Config struct {
	APIPort     string
	ProjectName string
	DatabaseURL string
	DebugMode   bool
}

// LoadConfig loads the configuration values from environment variables or the .env file
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	apiPort := getEnv("API_PORT", "8080")
	projectName := getEnv("PROJECT_NAME", "")
	databaseURL := getEnv("DATABASE_URL", "")
	debugMode, err := strconv.ParseBool(getEnv("DEBUG_MODE", "false"))
	if err != nil {
		log.Println("Failed to parse DEBUG_MODE. Defaulting to false.")
		debugMode = false
	}

	return &Config{
		APIPort:     apiPort,
		ProjectName: projectName,
		DatabaseURL: databaseURL,
		DebugMode:   debugMode,
	}
}

// getEnv retrieves the value of an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
