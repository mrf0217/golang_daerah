package http

import (
	"io"
	"net/http"
	"strconv"

	"context"
	"encoding/json"
	"golang_daerah/config"

	"github.com/jmoiron/sqlx"
)

type MySQLTrafficTicketSQLXRepository struct {
	db *sqlx.DB
}

func NewMySQLTrafficTicketSQLXRepository(db *sqlx.DB) *MySQLTrafficTicketSQLXRepository {
	return &MySQLTrafficTicketSQLXRepository{db: db}
}

func (r *MySQLTrafficTicketSQLXRepository) GetPaginatedJSON(limit, offset int) ([]byte, error) {
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

func (r *MySQLTrafficTicketSQLXRepository) InsertJSON(jsonData []byte) error {
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

func (h *MySQLTrafficTicketSQLXRepository) GetPaginated(w http.ResponseWriter, r *http.Request) {
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
		WriteInternalServerError(w, "Failed to get tickets: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *MySQLTrafficTicketSQLXRepository) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteBadRequest(w, "Invalid request body")
		return
	}

	if err := h.InsertJSON(body); err != nil {
		WriteInternalServerError(w, "Failed to insert tickets: "+err.Error())
		return
	}

	WriteSuccessResponseCreated(w, []interface{}{}, "Tickets created successfully")
}
