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

type AppConfig struct {
	PostgresPassword string
	PostgresUser     string
	PostgresDB       string
	PostgresHost     string
	PostgresPort     int
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
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     postgresPort,
	}

	if err := validateDatabaseConfig(Config); err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		if os.Getenv("APP_ENV") == "production" {
			log.Fatal("JWT_SECRET must be set in production")
		}
		jwtSecret = "dev-insecure-secret"
		log.Println("config: JWT_SECRET unset — insecure default in use. Set JWT_SECRET in production.")
	}
	JWTKey = []byte(jwtSecret)

	log.Println("config: loaded")
}

func loadEnvFile() {
	// In production, environment variables are injected by the orchestrator.
	// Skip .env file discovery to avoid accidentally loading stale files.
	if os.Getenv("APP_ENV") == "production" {
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to determine current working directory: %v", err)
	}

	envPath, err := findEnvFile(cwd)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("config: no .env file found — environment variables used as-is")
			return
		}
		log.Fatalf("failed to locate .env file: %v", err)
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("failed to load .env file %q: %v", envPath, err)
	}

	log.Printf("config: env loaded from %s", envPath)
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

	if cfg.PostgresUser == "" {
		missing = append(missing, "POSTGRES_USER")
	}
	if cfg.PostgresPassword == "" {
		missing = append(missing, "POSTGRES_PASSWORD")
	}
	if cfg.PostgresDB == "" {
		missing = append(missing, "POSTGRES_DB")
	}
	if cfg.PostgresHost == "" {
		missing = append(missing, "POSTGRES_HOST")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required database config: %s", strings.Join(missing, ", "))
	}

	return nil
}
