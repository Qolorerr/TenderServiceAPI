package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"zadanie_6105/src/models"
	"zadanie_6105/src/routes"
	"zadanie_6105/src/services"
)

var db *gorm.DB

func initDB() {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USERNAME"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_PORT"))

	var err error
	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to the database with GORM")

	err = db.AutoMigrate(&models.Tender{})
	if err != nil {
		log.Fatalf("Failed to migrate table: %v", err)
	}
	err = db.AutoMigrate(&models.Bid{})
	if err != nil {
		log.Fatalf("Failed to migrate table: %v", err)
	}
	err = db.AutoMigrate(&models.BidFeedback{})
	if err != nil {
		log.Fatalf("Failed to migrate table: %v", err)
	}
}

func main() {
	initDB()

	service := services.NewService(db)
	router := routes.RegisterRoutes(service)

	// Запуск HTTP сервера
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
