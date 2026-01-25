package repositories

import (
	"context"

	"pilates-reservation-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CourtRepository struct {
	db *pgxpool.Pool
}

func NewCourtRepository(db *pgxpool.Pool) *CourtRepository {
	return &CourtRepository{db: db}
}

func (r *CourtRepository) GetAvailable(
	ctx context.Context,
	date string,
	timeslotID int,
) ([]domain.Court, error) {

	rows, err := r.db.Query(ctx, `
		SELECT c.id, c.name
		FROM courts c
		WHERE c.id NOT IN (
			SELECT r.court_id
			FROM reservations r
			WHERE r.reservation_date = $1
			  AND r.timeslot_id = $2
		)
		ORDER BY c.name
	`, date, timeslotID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courts []domain.Court
	for rows.Next() {
		var c domain.Court
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		courts = append(courts, c)
	}

	return courts, nil
}
