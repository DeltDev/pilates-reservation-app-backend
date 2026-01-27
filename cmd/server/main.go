package main

import (
	"pilates-reservation-backend/internal/config"
	"pilates-reservation-backend/internal/db"
	"pilates-reservation-backend/internal/router"
	
)

func main() {
	cfg := config.Load()

	database := db.Connect(cfg.DBSource)
	defer database.Close()

	r := router.Setup(database)

	r.Run(":8080")
}
