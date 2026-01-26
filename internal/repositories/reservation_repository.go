package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationRepository struct {
	db *pgxpool.Pool
}

func NewReservationRepository(db *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) Create(
	ctx context.Context,
	date string,
	timeslotID, courtID int,
	customerName, customerEmail string,
) (int, error) {

	var id int

	err := r.db.QueryRow(ctx, `
		INSERT INTO reservations (
			reservation_date,
			timeslot_id,
			court_id,
			customer_name,
			customer_email
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`,
		date,
		timeslotID,
		courtID,
		customerName,
		customerEmail,
	).Scan(&id)

	return id, err
}
