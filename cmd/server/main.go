package server
import (
	"log"

	"pilates-booking-backend/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}

	db.RunMigrations()

	r := gin.Default()
	r.Run(":8080")
}
