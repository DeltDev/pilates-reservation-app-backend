package handlers

import (
	"net/http"

	"pilates-reservation-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

type TimeslotHandler struct {
	repo *repositories.TimeslotRepository
}

func NewTimeslotHandler(repo *repositories.TimeslotRepository) *TimeslotHandler {
	return &TimeslotHandler{repo: repo}
}

func (h *TimeslotHandler) GetTimeslots(c *gin.Context) {

	_, exists := c.GetQuery("date")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "date query parameter is required",
		})
		return
	}

	timeslots, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch timeslots",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"timeslots": timeslots,
	})
}
