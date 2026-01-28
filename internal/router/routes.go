package router

import (
	"pilates-reservation-backend/internal/handlers"
	"pilates-reservation-backend/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gin-contrib/cors"
)

func Setup(db *pgxpool.Pool) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
    	AllowOrigins: []string{
        	"http://localhost:3000",
			"https://zen-pilates.vercel.app", 
    	},
    	AllowMethods: []string{"GET", "POST", "OPTIONS"},
    	AllowHeaders: []string{"Content-Type"},
	}))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	timeslotRepo := repositories.NewTimeslotRepository(db)
	timeslotHandler := handlers.NewTimeslotHandler(timeslotRepo)

	courtRepo := repositories.NewCourtRepository(db)
	courtHandler := handlers.NewCourtHandler(courtRepo)

	reservationRepo := repositories.NewReservationRepository(db)
	reservationHandler := handlers.NewReservationHandler(reservationRepo)
	api := r.Group("/api")
	{
		api.GET("/timeslots", timeslotHandler.GetTimeslots)
		api.GET("/courts/available", courtHandler.GetAvailableCourts)
		api.POST("/reservations", reservationHandler.Create)
		api.GET("/courts", courtHandler.GetAllCourts)
		api.GET("/courts/:id", courtHandler.GetCourtAvailability)
	}

	return r
}
