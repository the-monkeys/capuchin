package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// AppConfig holds the configuration values for the application,

type AppConfig struct {
	POSTGRES_PASSWORD string
	POSTGRES_USER     string
	POSTGRES_DB       string
	POSTGRES_HOST     string
	JWTKey            string
}

var Config AppConfig
var JWTKey []byte

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Config = AppConfig{
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),
		POSTGRES_HOST:     os.Getenv("POSTGRES_HOST"),
	}

	// Load JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret" // Default for dev if env not set
		log.Println("WARNING: JWT_SECRET not set or empty; using default insecure secret. Set JWT_SECRET in production.")
	}
	Config.JWTKey = jwtSecret
	JWTKey = []byte(jwtSecret)

	log.Println("Configuration loaded successfully.")
}
