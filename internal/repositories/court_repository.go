package repositories

import (
	"context"
	"pilates-reservation-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)
type Timeslot struct {
	ID        int    `json:"id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type CourtWithAvailability struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
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

func (r *CourtRepository) GetAll() ([]domain.Court, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT id, name
		FROM courts
		ORDER BY id
	`)
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

func (r *CourtRepository) FindAvailableTimeslots(
	ctx context.Context,
	courtID int,
	date string,
) (CourtWithAvailability, []Timeslot, error) {

	var court CourtWithAvailability

	err := r.db.QueryRow(ctx, `
		SELECT id, name
		FROM courts
		WHERE id = $1
	`, courtID).Scan(&court.ID, &court.Name)

	if err != nil {
		return court, nil, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT t.id, t.start_time, t.end_time
		FROM timeslots t
		WHERE t.id NOT IN (
			SELECT timeslot_id
			FROM reservations
			WHERE reservation_date = $1
			  AND court_id = $2
		)
		ORDER BY t.start_time
	`, date, courtID)

	if err != nil {
		return court, nil, err
	}
	defer rows.Close()

	var timeslots []Timeslot
	for rows.Next() {
		var t Timeslot
		if err := rows.Scan(&t.ID, &t.StartTime, &t.EndTime); err != nil {
			return court, nil, err
		}
		timeslots = append(timeslots, t)
	}

	return court, timeslots, nil
}