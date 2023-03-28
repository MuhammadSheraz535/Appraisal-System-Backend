package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func Connect() {

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Create connection string using environment variables
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s database=postgres sslmode=disable", dbHost, dbUser, dbPassword, dbPort)

	// Open database connection
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Check if the database exists
	var count int64
	err = DB.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", dbName).Count(&count).Error
	if err != nil {
		panic(err)
	}

	// Create a new database if it doesn't exist
	if count == 0 {
		err = DB.Exec("CREATE DATABASE " + dbName).Error
		if err != nil {
			panic(err)
		}
	}

	// Close the current connection before connecting to the new database
	db, err := DB.DB()
	if err != nil {
		panic(err)
	}
	db.Close()

	// Use the new database
	dsn = fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s sslmode=disable", dbHost, dbUser, dbPassword, dbPort, dbName)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}
