package main

import (
	"fmt"
	"go-plate/models"
	"go-plate/routing"
	"go-plate/services"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	services.InitLogger()

	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "sqlite"
	}
	db, err := services.NewDB(&services.DatabaseConfig{
		Driver:   driver,
		Host:     os.Getenv("DB_HOST"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
	})

	if err != nil {
		fmt.Println("Database connexion failed", err.Error())
	} else {
		if db == nil {
			fmt.Println("Database is nil")
		} else {
			if migrate := models.MigrateModels(); !migrate {
				panic("Models migration failed")
			}
		}
	}

	router := mux.NewRouter()
	routing.RegisterRoutes(router)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	services.Logger.Println("Server is starting on port", port)
	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}