package handlers

import (
	"net/http"

	"pilates-reservation-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

type ReservationHandler struct {
	repo *repositories.ReservationRepository
}

func NewReservationHandler(repo *repositories.ReservationRepository) *ReservationHandler {
	return &ReservationHandler{repo: repo}
}

type createReservationRequest struct {
	Date         string `json:"date"`
	TimeslotID   int    `json:"timeslot_id"`
	CourtID      int    `json:"court_id"`
	CustomerName string `json:"customer_name"`
}

func (h *ReservationHandler) Create(c *gin.Context) {
	var req createReservationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if req.Date == "" || req.CustomerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing required fields",
		})
		return
	}

	err := h.repo.Create(
		c.Request.Context(),
		req.Date,
		req.TimeslotID,
		req.CourtID,
		req.CustomerName,
	)

	if err == repositories.ErrSlotAlreadyBooked {
		c.JSON(http.StatusConflict, gin.H{
			"error": "court already booked for this timeslot",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create reservation",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "reservation created successfully",
	})
}
