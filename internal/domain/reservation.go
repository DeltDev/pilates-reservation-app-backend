package domain

import "time"

type Reservation struct {
	ID              int       `json:"id"`
	ReservationDate string    `json:"reservation_date"`
	TimeslotID      int       `json:"timeslot_id"`
	CourtID         int       `json:"court_id"`
	CustomerName    string    `json:"customer_name"`
	CreatedAt       time.Time `json:"created_at"`
}
