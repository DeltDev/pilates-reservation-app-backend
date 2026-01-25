package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrSlotAlreadyBooked = errors.New("court already booked")

type ReservationRepository struct {
	db *pgxpool.Pool
}

func NewReservationRepository(db *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) Create(
	ctx context.Context,
	date string,
	timeslotID int,
	courtID int,
	customerName string,
) error {

	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM reservations
			WHERE reservation_date = $1
			  AND timeslot_id = $2
			  AND court_id = $3
		)
	`, date, timeslotID, courtID).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		return ErrSlotAlreadyBooked
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO reservations (
			reservation_date,
			timeslot_id,
			court_id,
			customer_name
		)
		VALUES ($1, $2, $3, $4)
	`, date, timeslotID, courtID, customerName)

	return err
}
