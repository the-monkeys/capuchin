package config

import "os"

type AppConfig struct {
	POSTGRES_PASSWORD string
	POSTGRES_USER     string
	POSTGRES_DB       string
	POSTGRES_HOST     string
}

var Config AppConfig

func init() {
	Config = AppConfig{
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),
		POSTGRES_HOST:     os.Getenv("POSTGRES_HOST"),
	}
}
