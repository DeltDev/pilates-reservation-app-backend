package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"

	"pilates-reservation-backend/internal/repositories"
)

type ReservationHandler struct {
	repo *repositories.ReservationRepository
}

func NewReservationHandler(repo *repositories.ReservationRepository) *ReservationHandler {
	return &ReservationHandler{repo: repo}
}

type CreateReservationRequest struct {
	Date          string `json:"date" binding:"required"`
	TimeslotID    int    `json:"timeslot_id" binding:"required"`
	CourtID       int    `json:"court_id" binding:"required"`
	CustomerName  string `json:"customer_name" binding:"required"`
	CustomerEmail string `json:"customer_email" binding:"required,email"`
}

func (h *ReservationHandler) Create(c *gin.Context) {
	var req CreateReservationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	id, err := h.repo.Create(
		c.Request.Context(),
		req.Date,
		req.TimeslotID,
		req.CourtID,
		req.CustomerName,
		req.CustomerEmail,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{
				"error": "court already reserved",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create reservation",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"reservation_id": id,
		"status":         "confirmed",
	})
}
