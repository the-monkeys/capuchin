package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// AppConfig holds the configuration values for the application,

type AppConfig struct {
	POSTGRES_PASSWORD string
	POSTGRES_USER     string
	POSTGRES_DB       string
	POSTGRES_HOST     string
	POSTGRES_PORT     int
	JWTKey            string
}

var Config AppConfig
var JWTKey []byte

func init() {
	loadEnvFile()

	postgresPort, err := postgresPortFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	Config = AppConfig{
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),
		POSTGRES_HOST:     os.Getenv("POSTGRES_HOST"),
		POSTGRES_PORT:     postgresPort,
	}

	if err := validateDatabaseConfig(Config); err != nil {
		log.Fatal(err)
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

func loadEnvFile() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to determine current working directory: %v", err)
	}

	envPath, err := findEnvFile(cwd)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("No .env file found in current or parent directories; using existing environment variables.")
			return
		}
		log.Fatalf("failed to locate .env file: %v", err)
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("failed to load .env file %q: %v", envPath, err)
	}

	log.Printf("Loaded environment variables from %s", envPath)
}

func findEnvFile(startDir string) (string, error) {
	dir := startDir
	for {
		candidate := filepath.Join(dir, ".env")
		info, err := os.Stat(candidate)
		if err == nil && !info.IsDir() {
			return candidate, nil
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", os.ErrNotExist
}

func postgresPortFromEnv() (int, error) {
	rawPort := os.Getenv("POSTGRES_PORT")
	if rawPort == "" {
		return 5432, nil
	}

	port, err := strconv.Atoi(rawPort)
	if err != nil {
		return 0, fmt.Errorf("invalid POSTGRES_PORT %q: must be numeric", rawPort)
	}
	if port <= 0 {
		return 0, fmt.Errorf("invalid POSTGRES_PORT %q: must be greater than 0", rawPort)
	}

	return port, nil
}

func validateDatabaseConfig(cfg AppConfig) error {
	missing := make([]string, 0, 4)

	if cfg.POSTGRES_USER == "" {
		missing = append(missing, "POSTGRES_USER")
	}
	if cfg.POSTGRES_PASSWORD == "" {
		missing = append(missing, "POSTGRES_PASSWORD")
	}
	if cfg.POSTGRES_DB == "" {
		missing = append(missing, "POSTGRES_DB")
	}
	if cfg.POSTGRES_HOST == "" {
		missing = append(missing, "POSTGRES_HOST")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required database config: %s", strings.Join(missing, ", "))
	}

	return nil
}
