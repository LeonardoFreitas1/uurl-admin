package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/LeonardoFreitas1/uurl-admin/db/sqlc"
	"github.com/joho/godotenv"
)

var (
	db      *sql.DB
	Queries *sqlc.Queries
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, proceeding with system environment variables")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbName == "" {
		log.Fatal("Database configuration variables are missing")
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbName,
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	Queries = sqlc.New(db)
}

func GetDB() *sql.DB {
	return db
}

func GetQueries() *sqlc.Queries {
	return Queries
}
