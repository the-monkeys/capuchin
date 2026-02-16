package db

import (
	"capuchin/config"
	"database/sql"
	"fmt"
	"time"
)

func NewConnection() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s", config.Config.POSTGRES_HOST, config.Config.POSTGRES_DB, config.Config.POSTGRES_USER, config.Config.POSTGRES_PASSWORD)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
