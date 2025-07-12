package main

import (
	"database/sql"
	"log"
	"rest-api-in-gin/internal/database"
	"rest-api-in-gin/internal/env"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

type application struct {
	port int
	jwtSecret string
	models database.Models
	DB        *sql.DB  
}

func main() {
	// fmt.Println("Hello")

	dsn := "host=localhost user=postgres password=postgres dbname=makerble port=5432 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		port: env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123456"),
		models: models,
		DB:        db,
	}


	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
	
}