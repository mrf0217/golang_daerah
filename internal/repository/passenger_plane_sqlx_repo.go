package repository

import (
    "context"
    "encoding/json"
    "golang_daerah/config"
    "github.com/jmoiron/sqlx"
)

type PassengerPlaneSQLXRepository struct {
    db *sqlx.DB
}

func NewPassengerPlaneSQLXRepository(db *sqlx.DB) *PassengerPlaneSQLXRepository {
    return &PassengerPlaneSQLXRepository{db: db}
}

func (r *PassengerPlaneSQLXRepository) GetPaginatedJSON(limit, offset int) ([]byte, error) {
    ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
    defer cancel()

    query := `
        SELECT id, passenger_name, passenger_id, age, gender, passport_number, 
               nationality, flight_number, departure_airport, arrival_airport, 
               departure_date, departure_time, arrival_time, seat_number, 
               ticket_class, baggage_weight, airline, gate, boarding_status,
               officer_name, officer_id, officer_rank, officer_branch_office_address, 
               checkin_counter, special_request
        FROM passenger_plane
        ORDER BY id ASC
        LIMIT ? OFFSET ?
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

func (r *PassengerPlaneSQLXRepository) InsertJSON(jsonData []byte) error {
    ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
    defer cancel()

    var items []map[string]interface{}
    if err := json.Unmarshal(jsonData, &items); err != nil {
        return err
    }

    query := `
        INSERT INTO passenger_plane (
            passenger_name, passenger_id, age, gender, passport_number, nationality,
            flight_number, departure_airport, arrival_airport, departure_date, 
            departure_time, arrival_time, seat_number, ticket_class, baggage_weight, 
            airline, gate, boarding_status, officer_name, officer_id, officer_rank, 
            officer_branch_office_address, checkin_counter, special_request
        ) VALUES (
            :passenger_name, :passenger_id, :age, :gender, :passport_number, :nationality,
            :flight_number, :departure_airport, :arrival_airport, :departure_date,
            :departure_time, :arrival_time, :seat_number, :ticket_class, :baggage_weight,
            :airline, :gate, :boarding_status, :officer_name, :officer_id, :officer_rank,
            :officer_branch_office_address, :checkin_counter, :special_request
        )
    `

    for _, item := range items {
        if _, err := r.db.NamedExecContext(ctx, query, item); err != nil {
            return handleQueryError(err)
        }
    }

    return nil
}