package router

import (
	"pilates-reservation-backend/internal/config"
	"pilates-reservation-backend/internal/handlers"
	"pilates-reservation-backend/internal/repositories"
	"pilates-reservation-backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(db *pgxpool.Pool, cfg config.Config) *gin.Engine {
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

	paymentService := services.NewPaymentService(
		cfg.MidtransServerKey,
		cfg.MidtransEnvironment,
	)

	timeslotRepo := repositories.NewTimeslotRepository(db)
	timeslotHandler := handlers.NewTimeslotHandler(timeslotRepo)

	courtRepo := repositories.NewCourtRepository(db)
	courtHandler := handlers.NewCourtHandler(courtRepo)

	reservationRepo := repositories.NewReservationRepository(db)
	reservationHandler := handlers.NewReservationHandler(reservationRepo, paymentService)

	api := r.Group("/api")
	{
		api.GET("/timeslots", timeslotHandler.GetTimeslots)
		api.GET("/courts/available", courtHandler.GetAvailableCourts)
		api.GET("/courts", courtHandler.GetAllCourts)
		api.GET("/courts/:id", courtHandler.GetCourtAvailability)
		api.POST("/reservations", reservationHandler.Create)
		api.POST("/payment/notification", reservationHandler.HandlePaymentNotification)
		api.GET("/payment/status", reservationHandler.CheckPaymentStatus)
	}

	r.GET("/api/config/midtrans", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"client_key":  cfg.MidtransClientKey,
			"environment": cfg.MidtransEnvironment,
		})
	})

	return r
}