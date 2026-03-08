package database

import (
<<<<<<< HEAD
	"capuchin/config"
	"database/sql"
	"fmt"
	"log"
=======
	"database/sql"
	"fmt"
	"log"
	"os"
>>>>>>> origin/auth-local
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		config.Config.POSTGRES_HOST,
		config.Config.POSTGRES_USER,
		config.Config.POSTGRES_PASSWORD,
		config.Config.POSTGRES_DB,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)
}

func InitSchema() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS todos (
		id UUID PRIMARY KEY,
		item TEXT NOT NULL,
		completed BOOLEAN DEFAULT FALSE,
		user_id UUID REFERENCES users(id)
	);`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Failed to init db:", err)
	}
}
