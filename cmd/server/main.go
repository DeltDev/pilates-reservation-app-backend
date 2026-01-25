package main

import (
	"log"

	"pilates-reservation-backend/internal/config"
	"pilates-reservation-backend/internal/db"
	"pilates-reservation-backend/internal/router"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	r := router.Setup(database)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
