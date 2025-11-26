-- name: InsertPassenger :execresult
INSERT INTO passenger_plane (
    passenger_name, passenger_id, age, gender, passport_number, nationality,
    flight_number, departure_airport, arrival_airport, departure_date, departure_time,
    arrival_time, seat_number, ticket_class, baggage_weight, airline, gate, boarding_status,
    officer_name, officer_id, officer_rank, officer_branch_office_address, checkin_counter, special_request
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?
);

-- name: ListPassengers :many
SELECT * FROM passenger_plane
ORDER BY id DESC
LIMIT ? OFFSET ?;

-- name: GetPaginatedPassengers :many
SELECT * FROM passenger_plane
ORDER BY id ASC
LIMIT ? OFFSET ?;

