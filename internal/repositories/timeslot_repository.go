package repositories

import (
	"context"

	"pilates-reservation-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TimeslotRepository struct {
	db *pgxpool.Pool
}

func NewTimeslotRepository(db *pgxpool.Pool) *TimeslotRepository {
	return &TimeslotRepository{db: db}
}

func (r *TimeslotRepository) GetAll(ctx context.Context) ([]domain.Timeslot, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, start_time, end_time
		FROM timeslots
		ORDER BY start_time
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timeslots []domain.Timeslot
	for rows.Next() {
		var t domain.Timeslot
		if err := rows.Scan(&t.ID, &t.StartTime, &t.EndTime); err != nil {
			return nil, err
		}
		timeslots = append(timeslots, t)
	}

	return timeslots, nil
}
