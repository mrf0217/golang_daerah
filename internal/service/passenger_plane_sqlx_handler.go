package service

import (
	"encoding/json"
	"golang_daerah/internal/database"
)

type PassengerPlaneService struct {
	db *database.BaseMultiDBRepository
}

func NewPassengerPlaneService(db *database.BaseMultiDBRepository) *PassengerPlaneService {
	return &PassengerPlaneService{db: db}
}

// func initializeDatabasesPassengerSQL() map[string]*sqlx.DB {
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

// func NewPassengerPlaneSQLXRepository(db *sqlx.DB) *PassengerPlaneSQLXRepository {
// 	return &PassengerPlaneSQLXRepository{db: db}
// }

func (r *PassengerPlaneService) GetPaginated(limit, offset int) ([]map[string]interface{}, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	// defer cancel()
	// db := r.getDB(dbName)
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

	result, err := r.db.QueryDB("passenger", query, limit, offset)
	if err != nil {
		return nil, err
	}

	// for i, port := range result {
	// 	// Database 2: Passengers
	// 	passengers, _ := r.db.QueryDB("passenger",
	// 		`SELECT passenger_name FROM passenger_plane WHERE id = ?`,
	// 		port["id"])

	// 	// Database 3: Traffic tickets
	// 	tickets, _ := r.db.QueryDB("traffic",
	// 		`SELECT legal_speed FROM traffic_tickets WHERE id = ?`,
	// 		port["id"])

	// 	// Database 4: Auth/Users (if needed)
	// 	users, _ := r.db.QueryDB("golang",
	// 		`SELECT username FROM users WHERE id = ?`,
	// 		port["id"])

	// 	result[i]["passengers"] = passengers
	// 	result[i]["traffic_ticket"] = tickets
	// 	result[i]["golang"] = users
	// }

	return result, nil
	// defer rows.Close()

	// var results []map[string]interface{}
	// for rows.Next() {
	// 	row := make(map[string]interface{})
	// 	if err := rows.MapScan(row); err != nil {
	// 		return nil, err
	// 	}
	// 	for key, value := range row {
	// 		if b, ok := value.([]byte); ok {
	// 			row[key] = string(b)
	// 		}
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

func (r *PassengerPlaneService) Create(jsonData []byte) error {
	// ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	// defer cancel()
	// db := r.getDB(dbName)
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
		if err := r.db.InsertDB("passenger", query, item); err != nil {
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

// func (h *PassengerPlaneSQLXRepository) GetPaginated_Passenger_SQL(w http.ResponseWriter, r *http.Request) {
// 	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
// 	perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))
// 	if page <= 0 {
// 		page = 1
// 	}
// 	if perPage <= 0 {
// 		perPage = 10
// 	}

// 	offset := (page - 1) * perPage
// 	jsonData, err := h.GetPaginatedJSON_Passenger_SQL(perPage, offset, "default")
// 	if err != nil {
// 		WriteInternalServerError(w, "Failed to get passengers: "+err.Error())
// 		return
// 	}

// 	var data []map[string]interface{}
// 	if err := json.Unmarshal(jsonData, &data); err != nil {
// 		WriteInternalServerError(w, "Failed to parse response: "+err.Error())
// 		return
// 	}

// 	WritePaginatedResponse(w, data, page, perPage, "Complete data retrieved successfully")
// }

// func (h *PassengerPlaneSQLXRepository) Create_Passenger_SQL(w http.ResponseWriter, r *http.Request) {
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		WriteBadRequest(w, "Invalid request body")
// 		return
// 	}

// 	if err := h.InsertJSON_Passenger_SQL(body, "default"); err != nil {
// 		WriteInternalServerError(w, "Failed to insert passengers: "+err.Error())
// 		return
// 	}

// 	WriteSuccessResponseCreated(w, []interface{}{}, "Passengers created successfully")
// }
