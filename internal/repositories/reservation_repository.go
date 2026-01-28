package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationRepository struct {
	db *pgxpool.Pool
}

func NewReservationRepository(db *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{db: db}
}

type CreateReservationParams struct {
	Date          string
	TimeslotID    int
	CourtID       int
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	OrderID       string
	GrossAmount   float64
}

func (r *ReservationRepository) Create(
	ctx context.Context,
	params CreateReservationParams,
) (int, error) {
	var id int

	err := r.db.QueryRow(ctx, `
		INSERT INTO reservations (
			reservation_date,
			timeslot_id,
			court_id,
			customer_name,
			customer_email,
			customer_phone,
			order_id,
			gross_amount,
			status,
			payment_status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending', 'pending')
		RETURNING id
	`,
		params.Date,
		params.TimeslotID,
		params.CourtID,
		params.CustomerName,
		params.CustomerEmail,
		params.CustomerPhone,
		params.OrderID,
		params.GrossAmount,
	).Scan(&id)

	return id, err
}

func (r *ReservationRepository) UpdatePaymentInfo(
	ctx context.Context,
	id int,
	token, redirectURL string,
) error {
	_, err := r.db.Exec(ctx, `
		UPDATE reservations 
		SET payment_token = $1, 
		    payment_redirect_url = $2,
		    updated_at = $3
		WHERE id = $4
	`, token, redirectURL, time.Now(), id)

	return err
}

func (r *ReservationRepository) UpdatePaymentStatus(
	ctx context.Context,
	orderID, paymentStatus, reservationStatus string,
) error {
	_, err := r.db.Exec(ctx, `
		UPDATE reservations 
		SET payment_status = $1,
		    status = $2,
		    updated_at = $3
		WHERE order_id = $4
	`, paymentStatus, reservationStatus, time.Now(), orderID)

	return err
}

func (r *ReservationRepository) GetByOrderID(
	ctx context.Context,
	orderID string,
) (*ReservationDetails, error) {
	var res ReservationDetails

	err := r.db.QueryRow(ctx, `
		SELECT 
			id, order_id, customer_name, customer_email, 
			payment_status, status, gross_amount
		FROM reservations
		WHERE order_id = $1
	`, orderID).Scan(
		&res.ID,
		&res.OrderID,
		&res.CustomerName,
		&res.CustomerEmail,
		&res.PaymentStatus,
		&res.Status,
		&res.GrossAmount,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

type ReservationDetails struct {
	ID            int
	OrderID       string
	CustomerName  string
	CustomerEmail string
	PaymentStatus string
	Status        string
	GrossAmount   float64
}