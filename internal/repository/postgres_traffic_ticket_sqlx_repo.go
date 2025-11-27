package repository

import (
    "context"
    "encoding/json"
    "golang_daerah/config"
    "github.com/jmoiron/sqlx"
)

type PostgresTrafficTicketSQLXRepository struct {
    db *sqlx.DB
}

func NewPostgresTrafficTicketSQLXRepository(db *sqlx.DB) *PostgresTrafficTicketSQLXRepository {
    return &PostgresTrafficTicketSQLXRepository{db: db}
}

func (r *PostgresTrafficTicketSQLXRepository) GetPaginatedJSON(limit, offset int) ([]byte, error) {
    ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
    defer cancel()

    query := `
        SELECT id, detected_speed, legal_speed, violation_location, 
               violation_date, violation_time, violation_type, 
               license_plate_number, vehicle_production_id, vehicle_factory,
               vehicle_model, vehicle_color, vehicle_brand, officer_name,
               officer_id, officer_rank, suspect_name, suspect_id, 
               suspect_age, officer_age, suspect_job, suspect_address,
               suspect_birth_place, officer_branch_office_address
        FROM traffic_tickets
        ORDER BY id ASC
        LIMIT $1 OFFSET $2
    `

    rows, err := r.db.QueryxContext(ctx, query, limit, offset)
    if err != nil {
        return nil, handleQueryError(err)
    }
    defer rows.Close()

    var results []map[string]interface{}
    for rows.Next() {
        row := make(map[string]interface{})
        if err := rows.MapScan(row); err != nil {
            return nil, err
        }
        results = append(results, row)
    }

    return json.Marshal(results)
}

func (r *PostgresTrafficTicketSQLXRepository) InsertJSON(jsonData []byte) error {
    ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
    defer cancel()

    var items []map[string]interface{}
    if err := json.Unmarshal(jsonData, &items); err != nil {
        return err
    }

    query := `
        INSERT INTO traffic_tickets (
            detected_speed, legal_speed, violation_location, violation_date, 
            violation_time, violation_type, license_plate_number, 
            vehicle_production_id, vehicle_factory, vehicle_model, vehicle_color,
            vehicle_brand, officer_name, officer_id, officer_rank, suspect_name,
            suspect_id, suspect_age, officer_age, suspect_job, suspect_address,
            suspect_birth_place, officer_branch_office_address
        ) VALUES (
            :detected_speed, :legal_speed, :violation_location, :violation_date,
            :violation_time, :violation_type, :license_plate_number,
            :vehicle_production_id, :vehicle_factory, :vehicle_model, :vehicle_color,
            :vehicle_brand, :officer_name, :officer_id, :officer_rank, :suspect_name,
            :suspect_id, :suspect_age, :officer_age, :suspect_job, :suspect_address,
            :suspect_birth_place, :officer_branch_office_address
        )
    `

    for _, item := range items {
        if _, err := r.db.NamedExecContext(ctx, query, item); err != nil {
            return handleQueryError(err)
        }
    }

    return nil
}