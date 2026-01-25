package handlers

import (
	"net/http"
	"strconv"

	"pilates-reservation-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

type CourtHandler struct {
	repo *repositories.CourtRepository
}

func NewCourtHandler(repo *repositories.CourtRepository) *CourtHandler {
	return &CourtHandler{repo: repo}
}

func (h *CourtHandler) GetAvailableCourts(c *gin.Context) {
	date := c.Query("date")
	timeslotIDStr := c.Query("timeslot_id")

	if date == "" || timeslotIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "date and timeslot_id are required",
		})
		return
	}

	timeslotID, err := strconv.Atoi(timeslotIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid timeslot_id",
		})
		return
	}

	courts, err := h.repo.GetAvailable(
		c.Request.Context(),
		date,
		timeslotID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch available courts",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courts": courts,
	})
}
