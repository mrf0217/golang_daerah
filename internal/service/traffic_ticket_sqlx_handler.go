package service

import (
	"encoding/json"
	"golang_daerah/internal/database"
)

type TrafficService struct {
	db *database.BaseMultiDBRepository
}

func NewPostgresTrafficTicketSQLXRepository(db *database.BaseMultiDBRepository) *TrafficService {
	return &TrafficService{db: db}
}

// func initializeDatabasesTrafficPostgre(db *database.BaseMultiDBRepository) map[string]*sqlx.DB {
// 	dbs := make(map[string]*sqlx.DB)

// 	// Default database
// 	dbs["default"] = config.InitTerminalDBX()

// 	// Add databases using existing config functions
// 	dbs["passenger"] = config.InitPassengerPlaneDBX()
// 	dbs["auth"] = config.InitAuthDBX()
// 	dbs["traffic"] = config.InitTrafficDBX()
// 	dbs["golang"] = config.InitGolangDBX()
// 	dbs["mysql"] = config.InitMySQLDBX()
// 	dbs["passanger"] = config.InitMySQLDBX_passanger()

// 	// Want to add a new database? Just add it here:
// 	// dbs["newdb"] = config.InitYourNewDBX()

// 	return dbs
// }

func (r *TrafficService) GetPaginated(limit, offset int) ([]map[string]interface{}, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	// defer cancel()
	// db := r.getDB(dbName)
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

	result, err := r.db.QueryDB("passenger", query, limit, offset)
	if err != nil {
		return nil, err
	}

	for i, port := range result {
		// Database 2: Passengers
		passengers, _ := r.db.QueryDB("passenger",
			`SELECT passenger_name FROM passenger_plane WHERE id = ?`,
			port["id"])

		// Database 3: Traffic tickets
		tickets, _ := r.db.QueryDB("traffic",
			`SELECT legal_speed FROM traffic_tickets WHERE id = ?`,
			port["id"])

		// Database 4: Auth/Users (if needed)
		users, _ := r.db.QueryDB("golang",
			`SELECT username FROM users WHERE id = ?`,
			port["id"])

		result[i]["passengers"] = passengers
		result[i]["traffic_ticket"] = tickets
		result[i]["golang"] = users
	}

	return result, nil
	// defer rows.Close()

	// var results []map[string]interface{}
	// for rows.Next() {
	// 	row := make(map[string]interface{})
	// 	if err := rows.MapScan(row); err != nil {
	// 		return nil, err
	// 	}
	// 	// portID := row["id"]

	// 	// // HARDCODED - ALWAYS queries these 3 databases
	// 	// passengers, _ := r.queryDB("passenger",
	// 	// 	`SELECT passenger_name FROM passenger_plane WHERE port_id = ?`,
	// 	// 	portID)

	// 	// tickets, _ := r.queryDB("traffic",
	// 	// 	`SELECT legal_speed FROM traffic_tickets WHERE port_id = ?`,
	// 	// 	portID)

	// 	// users, _ := r.queryDB("golang",
	// 	// 	`SELECT username FROM users WHERE port_id = ?`,
	// 	// 	portID)

	// 	// // ALWAYS adds these fields
	// 	// row["passengers"] = passengers
	// 	// row["tickets"] = tickets
	// 	// row["users"] = users

	// 	results = append(results, row)

	// }

	// return json.Marshal(results)
}

func (r *TrafficService) Create(jsonData []byte) error {
	// ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	// defer cancel()
	// db := r.getDB(dbName)
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
		if _, err := r.db.InsertDB("traffic", query, item); err != nil {
			return err
		}
	}

	// for _, item := range items {
	// 	// Insert into passenger database
	// 	passengerData := map[string]interface{}{
	// 		"passenger_name": item["suspect_name"],
	// 		"port_id":        item["id"],
	// 	}
	// 	r.insertDB("passenger",
	// 		`INSERT INTO passenger_plane (passenger_name, port_id) VALUES (:passenger_name, :port_id)`,
	// 		passengerData)

	// 	// Insert into traffic database
	// 	trafficData := map[string]interface{}{
	// 		"legal_speed": item["legal_speed"],
	// 		"port_id":     item["id"],
	// 	}
	// 	r.insertDB("traffic",
	// 		`INSERT INTO traffic_tickets (legal_speed, port_id) VALUES (:legal_speed, :port_id)`,
	// 		trafficData)

	// 	// Insert into golang/users database
	// 	userData := map[string]interface{}{
	// 		"username": item["officer_name"],
	// 		"port_id":  item["id"],
	// 	}
	// 	r.insertDB("golang",
	// 		`INSERT INTO users (username, port_id) VALUES (:username, :port_id)`,
	// 		userData)
	// }

	return nil
}

// func (h *PostgresTrafficTicketSQLXRepository) GetPaginated_Traffic_Postgre(w http.ResponseWriter, r *http.Request) {
// 	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
// 	perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))
// 	if page <= 0 {
// 		page = 1
// 	}
// 	if perPage <= 0 {
// 		perPage = 10
// 	}

// 	offset := (page - 1) * perPage
// 	jsonData, err := h.GetPaginatedJSON_Traffic_Postgre(perPage, offset, "default")
// 	if err != nil {
// 		WriteInternalServerError(w, "Failed to get tickets: "+err.Error())
// 		return
// 	}

// 	var data []map[string]interface{}
// 	if err := json.Unmarshal(jsonData, &data); err != nil {
// 		WriteInternalServerError(w, "Failed to parse response: "+err.Error())
// 		return
// 	}

// 	WritePaginatedResponse(w, data, page, perPage, "Complete data retrieved successfully")
// }

// func (h *PostgresTrafficTicketSQLXRepository) Create_Traffic_Postgre(w http.ResponseWriter, r *http.Request) {
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		WriteBadRequest(w, "Invalid request body")
// 		return
// 	}

// 	if err := h.InsertJSON_Traffic_Postgre(body, "default"); err != nil {
// 		WriteInternalServerError(w, "Failed to insert tickets: "+err.Error())
// 		return
// 	}

// 	WriteSuccessResponseCreated(w, []interface{}{}, "Tickets created successfully")
// }
