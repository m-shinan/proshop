package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func DbConnect() {
	dbURL := os.Getenv("DB_URL")
	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}

	// Verify the connection is valid
	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping DB: ", err)
	}

	log.Println("Successfully connected to the database")
}

// Automatically create tables if they don't exist

func CreateTables() error {
	// Read SQL file
	sqlFile, err := os.ReadFile(filepath.Join("database", "create_tables.sql"))
	if err != nil {
		return err
	}

	// Execute SQL commands
	_, err = DB.Exec(string(sqlFile))
	if err != nil {
		return err
	}

	log.Println("Tables created successfully")
	return nil
}
