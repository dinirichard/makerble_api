package main

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"rest-api-in-gin/internal/database"
	"rest-api-in-gin/internal/env"
	_ "github.com/lib/pq"
)

var testApp application

func TestMain(m *testing.M) {
	// DSN for the test database
	dsn := "host=localhost user=postgres password=postgres dbname=makerble_test port=5432 sslmode=disable"

	// Connect to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}

	// Set up the application instance for tests
	testApp = application{
		port:      env.GetEnvInt("PORT", 8081), // Use a different port for tests
		jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123456"),
		models:    database.NewModels(db),
		DB:        db,
	}

	// Run the tests
	code := m.Run()

	// Clean up the database after tests
	_, err = db.Exec("TRUNCATE TABLE patients, staffs RESTART IDENTITY CASCADE")
	if err != nil {
		log.Fatal("Failed to clean up test database:", err)
	}

	// Exit with the test result code
	os.Exit(code)
}
