package config

import (
	"log"
	"os"
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
		log.Println("config: JWT_SECRET unset - insecure default in use. Set JWT_SECRET in production.")
	}
	JWTKey = []byte(jwtSecret)

	log.Println("config: loaded")
}
