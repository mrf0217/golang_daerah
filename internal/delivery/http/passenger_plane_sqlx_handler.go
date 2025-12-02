package http

import (
	"context"
	"encoding/json"
	"golang_daerah/config"
	"io"
	"net/http"
	"strconv"

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
		for key, value := range row {
			if b, ok := value.([]byte); ok {
				row[key] = string(b)
			}
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

func (h *PassengerPlaneSQLXRepository) GetPaginated(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	jsonData, err := h.GetPaginatedJSON(perPage, offset)
	if err != nil {
		WriteInternalServerError(w, "Failed to get passengers: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *PassengerPlaneSQLXRepository) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteBadRequest(w, "Invalid request body")
		return
	}

	if err := h.InsertJSON(body); err != nil {
		WriteInternalServerError(w, "Failed to insert passengers: "+err.Error())
		return
	}

	WriteSuccessResponseCreated(w, []interface{}{}, "Passengers created successfully")
}
