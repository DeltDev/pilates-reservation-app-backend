package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Court struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CourtRepository struct {
	db *pgxpool.Pool
}

func NewCourtRepository(db *pgxpool.Pool) *CourtRepository {
	return &CourtRepository{db: db}
}

func (r *CourtRepository) FindAvailable(
	ctx context.Context,
	date string,
	timeslotID int,
) ([]Court, error) {

	rows, err := r.db.Query(ctx, `
		SELECT c.id, c.name
		FROM courts c
		WHERE c.id NOT IN (
			SELECT court_id
			FROM reservations
			WHERE reservation_date = $1 AND timeslot_id = $2
		)
		ORDER BY c.id
	`, date, timeslotID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courts []Court
	for rows.Next() {
		var c Court
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		courts = append(courts, c)
	}

	return courts, nil
}
