package infrastructure

import (
	"database/sql"
	"fmt"
	"os"
)

func SetupDatabaseConnection() *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", 
		os.Getenv("DB_HOST"), 
		os.Getenv("DB_USER"), 
		os.Getenv("DB_PASSWORD"), 
		os.Getenv("DB_NAME"), 
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
	// dsn := fmt.Sprintf(
	// 	"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", 
	// 	os.Getenv("DB_HOST"), 
	// 	os.Getenv("DB_USER"), 
	// 	os.Getenv("DB_PASSWORD"), 
	// 	os.Getenv("DB_NAME"), 
	// 	os.Getenv("DB_PORT"),
	// )

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		panic("Failed to create connection to database")
	}

	return db
}

func CloseDatabaseConnection(db *sql.DB) {
	db.Close()
}
