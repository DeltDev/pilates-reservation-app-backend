package main

import (
	"context"
	"log"
	"pilates-reservation-backend/internal/config"
	"pilates-reservation-backend/internal/db"
	"pilates-reservation-backend/internal/router"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	database := db.Connect(cfg.DBSource)
	defer database.Close()

	log.Println("Running auto-migration...")
	if err := initDatabase(database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Migration successful!")

	r := router.Setup(database, cfg)

	log.Println("Server starting on :8080")
	r.Run(":8080")
}

func initDatabase(db *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS courts (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS timeslots (
			id SERIAL PRIMARY KEY,
			start_time TIME NOT NULL,
			end_time TIME NOT NULL,
			CHECK (start_time < end_time)
		);`,

		`CREATE TABLE IF NOT EXISTS reservations (
			id SERIAL PRIMARY KEY,
			reservation_date DATE NOT NULL,
			timeslot_id INT NOT NULL,
			court_id INT NOT NULL,
			customer_name TEXT NOT NULL,
			customer_email TEXT NOT NULL,
			customer_phone TEXT,
			status TEXT NOT NULL DEFAULT 'pending',
			payment_status TEXT NOT NULL DEFAULT 'pending',
			payment_token TEXT,
			payment_redirect_url TEXT,
			order_id TEXT UNIQUE,
			gross_amount DECIMAL(10,2) NOT NULL DEFAULT 150000,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_timeslot FOREIGN KEY (timeslot_id) REFERENCES timeslots(id) ON DELETE RESTRICT,
			CONSTRAINT fk_court FOREIGN KEY (court_id) REFERENCES courts(id) ON DELETE RESTRICT,
			CONSTRAINT unique_reservation UNIQUE (reservation_date, timeslot_id, court_id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_order_id ON reservations(order_id);`,
		`CREATE INDEX IF NOT EXISTS idx_payment_status ON reservations(payment_status);`,
	}

	for _, q := range queries {
		_, err := db.Exec(ctx, q)
		if err != nil {
			return err
		}
	}

	var count int
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM courts").Scan(&count)
	if err == nil && count == 0 {
		log.Println("Seeding Courts...")
		_, err = db.Exec(ctx, `INSERT INTO courts (name) VALUES 
			('Court 1'), ('Court 2'), ('Court 3'), ('Court 4'), ('Court 5'),
			('Court 6'), ('Court 7'), ('Court 8'), ('Court 9'), ('Court 10');`) 
		if err != nil {
			return err
		}
	}

	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM timeslots").Scan(&count)
	if err == nil && count == 0 {
		log.Println("Seeding Timeslots...")
		_, err = db.Exec(ctx, `INSERT INTO timeslots (start_time, end_time) VALUES
			('07:00', '08:00'), ('08:00', '09:00'), ('09:00', '10:00'), ('10:00', '11:00'),
			('11:00', '12:00'), ('12:00', '13:00'), ('13:00', '14:00'), ('14:00', '15:00'),
			('15:00', '16:00'), ('16:00', '17:00'), ('17:00', '18:00'), ('18:00', '19:00'),
			('19:00', '20:00'), ('20:00', '21:00');`)
		if err != nil {
			return err
		}
	}

	return nil
}