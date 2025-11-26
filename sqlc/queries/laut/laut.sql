-- name: InsertLaut :execresult
INSERT INTO Laut (
    port_name, port_code, port_address, city, province, country,
    operator_name, operator_contact, harbor_master_name, harbor_master_id,
    harbor_master_rank, harbor_master_office_address, number_of_piers,
    main_pier_length, max_ship_draft, max_ship_length,
    terminal_capacity_passenger, terminal_capacity_cargo, operational_hours,
    emergency_contact, security_office_name, security_officer_id,
    security_level, checkin_counter_count, special_facilities
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?
);

-- name: ListLauts :many
SELECT * FROM Laut
ORDER BY id DESC
LIMIT ? OFFSET ?;

-- name: GetPaginatedLauts :many
SELECT * FROM Laut
ORDER BY id ASC
LIMIT ? OFFSET ?;

