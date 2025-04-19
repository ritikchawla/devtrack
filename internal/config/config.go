// internal/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDriver string
	DBSource string
	APIPort  string
}

// LoadConfig loads configuration from environment variables or .env file
func LoadConfig(path string) (config Config, err error) {
	// Attempt to load .env file if it exists
	err = godotenv.Load(path + "/.env")
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
		// Don't return error if .env is missing, just log it
		err = nil
	}

	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "postgres" // Default driver
	}

	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		log.Fatal("DB_SOURCE environment variable is required")
	}

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "8080" // Default port
	}

	config = Config{
		DBDriver: dbDriver,
		DBSource: dbSource,
		APIPort:  apiPort,
	}
	return
}
