-- name: InsertTrafficTicket :execresult
INSERT INTO traffic_tickets (
    detected_speed, legal_speed, violation_location, violation_date, violation_time,
    violation_type, license_plate_number, vehicle_production_id, vehicle_factory, vehicle_model,
    vehicle_color, vehicle_brand, officer_name, officer_id, officer_rank, suspect_name,
    suspect_id, suspect_age, officer_age, suspect_job, suspect_address, suspect_birth_place,
    officer_branch_office_address
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
    ?, ?, ?
);

-- name: ListTrafficTickets :many
SELECT * FROM traffic_tickets
ORDER BY id DESC
LIMIT ? OFFSET ?;

-- name: GetPaginatedTrafficTickets :many
SELECT * FROM traffic_tickets
ORDER BY id ASC
LIMIT ? OFFSET ?;

