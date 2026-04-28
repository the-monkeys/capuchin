package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	POSTGRES_PASSWORD string
	POSTGRES_USER     string
	POSTGRES_DB       string
	POSTGRES_HOST     string
	POSTGRES_PORT     int
}

var Config AppConfig
var JWTKey []byte

func init() {
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
		log.Println("WARNING: JWT_SECRET not set or empty; using default insecure secret. Set JWT_SECRET in production.")
	}
	JWTKey = []byte(jwtSecret)

	log.Println("Configuration loaded successfully.")
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
