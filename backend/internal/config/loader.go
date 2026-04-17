package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

func loadEnvFile() {
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
			log.Println("config: no .env file found - environment variables used as-is")
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
