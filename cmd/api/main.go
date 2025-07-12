package main

import (
	"database/sql"
	"log"
	"makerble_api/internal/database"
	"makerble_api/internal/env"

	_ "makerble_api/docs"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

// @title Makerble rest api
// @version 1.0
// @description A rest API in Go using Gin framework
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**

type application struct {
	port int
	jwtSecret string
	models database.Models
	DB        *sql.DB  
}

func main() {

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