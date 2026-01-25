CREATE TABLE reservations (
    id SERIAL PRIMARY KEY,
    reservation_date DATE NOT NULL,
    timeslot_id INT NOT NULL,
    court_id INT NOT NULL,

    customer_name TEXT NOT NULL,
    customer_email TEXT NOT NULL,

    status TEXT NOT NULL DEFAULT 'confirmed',

    CONSTRAINT fk_timeslot
        FOREIGN KEY (timeslot_id)
        REFERENCES timeslots(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_court
        FOREIGN KEY (court_id)
        REFERENCES courts(id)
        ON DELETE RESTRICT,

    CONSTRAINT unique_reservation
        UNIQUE (reservation_date, timeslot_id, court_id)
);
