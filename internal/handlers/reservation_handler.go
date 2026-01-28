package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"

	"pilates-reservation-backend/internal/repositories"
	"pilates-reservation-backend/internal/services"
)

type ReservationHandler struct {
	repo           *repositories.ReservationRepository
	paymentService *services.PaymentService
}

func NewReservationHandler(
	repo *repositories.ReservationRepository,
	paymentService *services.PaymentService,
) *ReservationHandler {
	return &ReservationHandler{
		repo:           repo,
		paymentService: paymentService,
	}
}

type CreateReservationRequest struct {
	Date          string `json:"date" binding:"required"`
	TimeslotID    int    `json:"timeslot_id" binding:"required"`
	CourtID       int    `json:"court_id" binding:"required"`
	CustomerName  string `json:"customer_name" binding:"required"`
	CustomerEmail string `json:"customer_email" binding:"required,email"`
	CustomerPhone string `json:"customer_phone" binding:"required"`
}

func (h *ReservationHandler) Create(c *gin.Context) {
	var req CreateReservationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	orderID := fmt.Sprintf("PILATES-%d", time.Now().UnixNano())
	grossAmount := 150000.0 

	reservationID, err := h.repo.Create(
		c.Request.Context(),
		repositories.CreateReservationParams{
			Date:          req.Date,
			TimeslotID:    req.TimeslotID,
			CourtID:       req.CourtID,
			CustomerName:  req.CustomerName,
			CustomerEmail: req.CustomerEmail,
			CustomerPhone: req.CustomerPhone,
			OrderID:       orderID,
			GrossAmount:   grossAmount,
		},
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{
				"error": "court already reserved for this timeslot",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create reservation",
		})
		return
	}

	paymentResp, err := h.paymentService.CreateTransaction(services.CreateTransactionRequest{
		OrderID:       orderID,
		GrossAmount:   int64(grossAmount),
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		CustomerPhone: req.CustomerPhone,
		ItemName:      "Pilates Session Booking",
		ItemPrice:     int64(grossAmount),
		ItemQuantity:  1,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create payment transaction",
		})
		return
	}

	err = h.repo.UpdatePaymentInfo(
		c.Request.Context(),
		reservationID,
		paymentResp.Token,
		paymentResp.RedirectURL,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update payment info",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"reservation_id":  reservationID,
		"order_id":        orderID,
		"payment_token":   paymentResp.Token,
		"payment_url":     paymentResp.RedirectURL,
		"status":          "pending",
		"payment_status":  "pending",
		"amount":          grossAmount,
	})
}

func (h *ReservationHandler) HandlePaymentNotification(c *gin.Context) {
	var notification map[string]interface{}

	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification"})
		return
	}

	orderID, _ := notification["order_id"].(string)
	transactionStatus, _ := notification["transaction_status"].(string)
	fraudStatus, _ := notification["fraud_status"].(string)

	var paymentStatus, reservationStatus string
	if transactionStatus == "capture" {
		if fraudStatus == "accept" {
			paymentStatus = "paid"
			reservationStatus = "confirmed"
		}
	} else if transactionStatus == "settlement" {
		paymentStatus = "paid"
		reservationStatus = "confirmed"
	} else if transactionStatus == "cancel" || transactionStatus == "deny" || transactionStatus == "expire" {
		paymentStatus = "failed"
		reservationStatus = "cancelled"
	} else if transactionStatus == "pending" {
		paymentStatus = "pending"
		reservationStatus = "pending"
	}
	err := h.repo.UpdatePaymentStatus(
		c.Request.Context(),
		orderID,
		paymentStatus,
		reservationStatus,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ReservationHandler) CheckPaymentStatus(c *gin.Context) {
	orderID := c.Query("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id is required"})
		return
	}

	reservation, err := h.repo.GetByOrderID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reservation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id":       reservation.OrderID,
		"payment_status": reservation.PaymentStatus,
		"status":         reservation.Status,
	})
}