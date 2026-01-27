package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pilates-reservation-backend/internal/repositories"
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

	courts, err := h.repo.FindAvailable(
		c.Request.Context(),
		date,
		timeslotID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courts": courts,
	})
}

func (h *CourtHandler) GetAllCourts(c *gin.Context) {
	courts, err := h.repo.GetAll()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to fetch courts",
		})
		return
	}

	c.JSON(200, gin.H{
		"courts": courts,
	})
}

func (h *CourtHandler) GetCourtAvailability(c *gin.Context) {
	courtIDStr := c.Param("id")
	date := c.Query("date")

	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "date is required",
		})
		return
	}

	courtID, err := strconv.Atoi(courtIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid court id",
		})
		return
	}

	court, timeslots, err := h.repo.FindAvailableTimeslots(
		c.Request.Context(),
		courtID,
		date,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "court not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"court":               court,
		"available_timeslots": timeslots,
	})
}