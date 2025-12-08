package service

import (
	"encoding/json"
	"golang_daerah/internal/database"
)

type LautService struct {
	db *database.BaseMultiDBRepository
}

// var joinQueries = map[string]string{
// 	"passenger": `SELECT passenger_name FROM passenger_plane WHERE id = ?`,
// 	"traffic":   `SELECT violation_location FROM traffic_tickets WHERE id = ?`,
// 	// "auth":      `SELECT username, role FROM Users WHERE port_id = ?`,
// 	// "golang":    `SELECT * FROM SomeTable WHERE port_id = ?`,
// }

// func NewLautRepository(db *sqlx.DB) *LautSQLXRepository {
// 	return &LautSQLXRepository{
// 		db:  db,
// 		dbs: LautinitializeDatabases(),
// 	}
// }

func NewLautService(db *database.BaseMultiDBRepository) *LautService {
	return &LautService{db: db}
}

// ADD YOUR DATABASES HERE - Just call the config functions!
// func LautinitializeDatabases() map[string]*sqlx.DB {
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

// -------------------------------------------------------------------------------------------
// Extract database name from URL path
// func extractDBName(path string) string {
// 	// Example: /api/terminals/passenger -> "passenger"
// 	// Example: /api/terminals -> "default"
// 	parts := strings.Split(strings.Trim(path, "/"), "/")

// 	if len(parts) >= 3 {
// 		return parts[2]
// 	}

// 	return "default"
// }
// -------------------------------------------------------------------------------------------
// func (h *LautSQLXRepository) List(w http.ResponseWriter, r *http.Request) {
//     jsonData, err := h.GetAllJSON()
//     if err != nil {
//         WriteInternalServerError(w, "Failed to get terminals: "+err.Error())
//         return
//     }

//	    w.Header().Set("Content-Type", "application/json")
//	    w.Write(jsonData)
//	}
//
// --------------------------------------------------------------------------------------------
func (r *LautService) Create(jsonData []byte) error {
	var items []map[string]interface{}
	if err := json.Unmarshal(jsonData, &items); err != nil {
		return err
	}

	query := `
        INSERT INTO Laut (
            port_name, port_code, port_address, city, province, country,
            operator_name, operator_contact, harbor_master_name, harbor_master_id,
            harbor_master_rank, harbor_master_office_address, number_of_piers,
            main_pier_length, max_ship_draft, max_ship_length,
            terminal_capacity_passenger, terminal_capacity_cargo, operational_hours,
            emergency_contact, security_office_name, security_officer_id,
            security_level, checkin_counter_count, special_facilities
        ) VALUES (
            :port_name, :port_code, :port_address, :city, :province, :country,
            :operator_name, :operator_contact, :harbor_master_name, :harbor_master_id,
            :harbor_master_rank, :harbor_master_office_address, :number_of_piers,
            :main_pier_length, :max_ship_draft, :max_ship_length,
            :terminal_capacity_passenger, :terminal_capacity_cargo, :operational_hours,
            :emergency_contact, :security_office_name, :security_officer_id,
            :security_level, :checkin_counter_count, :special_facilities
        )
    `

	for _, item := range items {
		if err := r.db.InsertDB("terminal", query, item); err != nil {
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

func (r *LautService) GetPaginated(limit, offset int) ([]map[string]interface{}, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	// defer cancel()
	// db := r.getDB(dbName)
	query := `
        SELECT id, port_name, port_code, port_address, city, province, country,
               operator_name, operator_contact, harbor_master_name, harbor_master_id,
               harbor_master_rank, harbor_master_office_address, number_of_piers,
               main_pier_length, max_ship_draft, max_ship_length,
               terminal_capacity_passenger, terminal_capacity_cargo, operational_hours,
               emergency_contact, security_office_name, security_officer_id,
               security_level, checkin_counter_count, special_facilities
        FROM Laut
        ORDER BY id ASC
        LIMIT ? OFFSET ?
    `
	result, err:= r.db.QueryDB("terminal", query, limit, offset)
	if err != nil {
        return nil, err
    }
	//this for multiple db query
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
	

	// rows, err := r.db.QueryxContext("terminal", query, limit, offset)
	// if err != nil {
	// 	return nil, handleQueryError(err)
	// }
	// defer rows.Close()

	// var results []map[string]interface{}
	// for rows.Next() {
	// 	row := make(map[string]interface{})
	// 	if err := rows.MapScan(row); err != nil {
	// 		return nil, err
	// 	}

	// 	// Convert []byte to string for all fields
	// 	for key, value := range row {
	// 		if b, ok := value.([]byte); ok {
	// 			row[key] = string(b)
	// 		}

	// 	}

	// portID := row["id"]

	// // HARDCODED - ALWAYS queries these 3 databases
	// passengers, _ := r.queryDB("passenger",
	// 	`SELECT passenger_name FROM passenger_plane WHERE id = ?`,
	// 	portID)

	// tickets, _ := r.queryDB("traffic",
	// 	`SELECT legal_speed FROM traffic_tickets WHERE id = ?`,
	// 	portID)

	// users, _ := r.queryDB("golang",
	// 	`SELECT username FROM users WHERE id = ?`,
	// 	portID)

	// // ALWAYS adds these fields
	// row["passengers"] = passengers
	// row["tickets"] = tickets
	// row["users"] = users
	// results = append(results, row)
	// }
	// return json.Marshal(results)
}

func (r *LautService) GetCompleteData(limit, offset int) ([]map[string]interface{}, error) {
	// Database 1: Ports
	ports, err := r.db.QueryDB("terminal",
		`SELECT id, port_name FROM Laut LIMIT ? OFFSET ?`,
		limit, offset)
	if err != nil {
		return nil, err
	}

	for i, port := range ports {
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

		ports[i]["passengers"] = passengers
		ports[i]["traffic_ticket"] = tickets
		ports[i]["golang"] = users
	}

	return ports, nil
}

// func (h *LautSQLXRepository) LautGetCompleteDataHandler(w http.ResponseWriter, r *http.Request) {
// 	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
// 	perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))
// 	if page <= 0 {
// 		page = 1
// 	}
// 	if perPage <= 0 {
// 		perPage = 10
// 	}

// 	offset := (page - 1) * perPage
// 	jsonData, err := h.LautGetCompleteData(perPage, offset)
// 	if err != nil {
// 		WriteInternalServerError(w, "Failed to get complete data: "+err.Error())
// 		return
// 	}

// 	var data []map[string]interface{}
// 	if err := json.Unmarshal(jsonData, &data); err != nil {
// 		WriteInternalServerError(w, "Failed to parse response: "+err.Error())
// 		return
// 	}

// 	WritePaginatedResponse(w, data, page, perPage, "Complete data retrieved successfully")
// }

// func (h *LautSQLXRepository) LautGetPaginated(w http.ResponseWriter, r *http.Request) {

// 	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
// 	perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))
// 	if page <= 0 {
// 		page = 1
// 	}
// 	if perPage <= 0 {
// 		perPage = 10
// 	}
// 	offset := (page - 1) * perPage

// 	jsonData, err := h.LautGetPaginatedJSON(perPage, offset, "default")
// 	if err != nil {
// 		WriteInternalServerError(w, "Failed to get terminals: "+err.Error())
// 		return
// 	}

// 	var data []map[string]interface{}
// 	if err := json.Unmarshal(jsonData, &data); err != nil {
// 		WriteInternalServerError(w, "Failed to parse response: "+err.Error())
// 		return
// 	}

// 	WritePaginatedResponse(w, data, page, perPage, "Complete data retrieved successfully")

// }

// func (h *LautSQLXRepository) Create(w http.ResponseWriter, r *http.Request) {
// 	// dbName := extractDBName(r.URL.Path)
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		WriteBadRequest(w, "Invalid request body")
// 		return
// 	}

// 	if err := h.LautInsertJSON(body, "default"); err != nil {
// 		WriteInternalServerError(w, "Failed to insert terminals: "+err.Error())
// 		return
// 	}

// 	WriteSuccessResponseCreated(w, []interface{}{}, "Terminals created successfully")
// }

// func New(db *sqlx.DB) *LautSQLXRepository {
// 	return &LautSQLXRepository{db: db}
// }

// func (r *LautSQLXRepository) GetAllJSON(limit, offset int) ([]byte, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
// 	defer cancel()

// 	query := `
//         SELECT id, port_name, port_code, port_address, city, province, country,
//                operator_name, operator_contact, harbor_master_name, harbor_master_id,
//                harbor_master_rank, harbor_master_office_address, number_of_piers,
//                main_pier_length, max_ship_draft, max_ship_length,
//                terminal_capacity_passenger, terminal_capacity_cargo, operational_hours,
//                emergency_contact, security_office_name, security_officer_id,
//                security_level, checkin_counter_count, special_facilities
//         FROM Laut
//         ORDER BY id ASC
//         LIMIT ? OFFSET ?
//     `

// 	rows, err := r.db.QueryxContext(ctx, query, limit, offset,"dab")
// 	if err != nil {
// 		return nil, handleQueryError(err)
// 	}
// 	defer rows.Close()

// 	var results []map[string]interface{}
// 	for rows.Next() {
// 		row := make(map[string]interface{})
// 		if err := rows.MapScan(row); err != nil {
// 			return nil, err
// 		}

// 		// Convert []byte to string for all fields
// 		for key, value := range row {
// 			if b, ok := value.([]byte); ok {
// 				row[key] = string(b)
// 			}

// 		}
// 		results = append(results, row)
// 	}
// 	return json.Marshal(results)
// }
