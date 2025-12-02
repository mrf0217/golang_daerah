package http

import (
	"context"
	"encoding/json"
	"fmt"
	"golang_daerah/config"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type LautSQLXRepository struct {
	db  *sqlx.DB
	dbs map[string]*sqlx.DB
}

// var joinQueries = map[string]string{
// 	"passenger": `SELECT passenger_name FROM passenger_plane WHERE id = ?`,
// 	"traffic":   `SELECT violation_location FROM traffic_tickets WHERE id = ?`,
// 	// "auth":      `SELECT username, role FROM Users WHERE port_id = ?`,
// 	// "golang":    `SELECT * FROM SomeTable WHERE port_id = ?`,
// }

func NewLautRepository(db *sqlx.DB) *LautSQLXRepository {
	return &LautSQLXRepository{
		db:  db,
		dbs: initializeDatabases(),
	}
}

// ADD YOUR DATABASES HERE - Just call the config functions!
func initializeDatabases() map[string]*sqlx.DB {
	dbs := make(map[string]*sqlx.DB)

	// Default database
	dbs["default"] = config.InitTerminalDBX()

	// Add databases using existing config functions
	dbs["passenger"] = config.InitPassengerPlaneDBX()
	dbs["auth"] = config.InitAuthDBX()
	dbs["traffic"] = config.InitTrafficDBX()
	dbs["golang"] = config.InitGolangDBX()
	dbs["mysql"] = config.InitMySQLDBX()
	dbs["passanger"] = config.InitMySQLDBX_passanger()

	// Want to add a new database? Just add it here:
	// dbs["newdb"] = config.InitYourNewDBX()

	return dbs
}

// Get database by name
func (r *LautSQLXRepository) getDB(dbName string) *sqlx.DB {
	if db, exists := r.dbs[dbName]; exists {
		return db
	}
	return r.dbs["default"]
}

// Extract database name from URL path
func extractDBName(path string) string {
	// Example: /api/terminals/passenger -> "passenger"
	// Example: /api/terminals -> "default"
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) >= 3 {
		return parts[2]
	}

	return "default"
}

func (h *LautSQLXRepository) GetCompleteDataHandler(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	jsonData, err := h.GetCompleteData(perPage, offset)
	if err != nil {
		WriteInternalServerError(w, "Failed to get complete data: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *LautSQLXRepository) GetPaginated(w http.ResponseWriter, r *http.Request) {
	// joinParam := r.URL.Query().Get("joinDbs")
	dbName := extractDBName(r.URL.Path)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	// "passenger,auth"

	// var joinDBs []string

	// if joinParam != "" {
	// 	joinDBs = strings.Split(joinParam, ",")
	// }

	jsonData, err := h.GetPaginatedJSON(perPage, offset, dbName)
	if err != nil {
		WriteInternalServerError(w, "Failed to get terminals: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *LautSQLXRepository) Create(w http.ResponseWriter, r *http.Request) {
	dbName := extractDBName(r.URL.Path)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteBadRequest(w, "Invalid request body")
		return
	}

	if err := h.InsertJSON(body, dbName); err != nil {
		WriteInternalServerError(w, "Failed to insert terminals: "+err.Error())
		return
	}

	WriteSuccessResponseCreated(w, []interface{}{}, "Terminals created successfully")
}

// func (h *LautSQLXRepository) List(w http.ResponseWriter, r *http.Request) {
//     jsonData, err := h.GetAllJSON()
//     if err != nil {
//         WriteInternalServerError(w, "Failed to get terminals: "+err.Error())
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     w.Write(jsonData)
// }

func (r *LautSQLXRepository) InsertJSON(jsonData []byte, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	db := r.getDB(dbName)

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
		if _, err := db.NamedExecContext(ctx, query, item); err != nil {
			return handleQueryError(err)
		}
	}

	return nil
}

func (r *LautSQLXRepository) GetPaginatedJSON(limit, offset int, dbName string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()
	db := r.getDB(dbName)
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

	rows, err := db.QueryxContext(ctx, query, limit, offset)
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

		// Convert []byte to string for all fields
		for key, value := range row {
			if b, ok := value.([]byte); ok {
				row[key] = string(b)
			}

		}

		portID := row["id"]

		// HARDCODED - ALWAYS queries these 3 databases
		passengers, _ := r.queryDB("passenger",
			`SELECT passenger_name FROM passenger_plane WHERE port_id = ?`,
			portID)

		tickets, _ := r.queryDB("traffic",
			`SELECT legal_speed FROM traffic_tickets WHERE port_id = ?`,
			portID)

		users, _ := r.queryDB("golang",
			`SELECT username FROM users WHERE port_id = ?`,
			portID)

		// ALWAYS adds these fields
		row["passengers"] = passengers
		row["tickets"] = tickets
		row["users"] = users
		results = append(results, row)
	}
	return json.Marshal(results)
}

func (r *LautSQLXRepository) GetCompleteData(limit, offset int) ([]byte, error) {
	// Database 1: Ports
	ports, err := r.queryDB("terminal",
		`SELECT id, port_name FROM Laut LIMIT ? OFFSET ?`,
		limit, offset)
	if err != nil {
		return nil, err
	}

	for i, port := range ports {
		// Database 2: Passengers
		passengers, _ := r.queryDB("passenger",
			`SELECT passenger_name FROM passenger_plane WHERE id = ?`,
			port["id"])

		// Database 3: Traffic tickets
		tickets, _ := r.queryDB("traffic",
			`SELECT legal_speed FROM traffic_tickets WHERE id = ?`,
			port["id"])

		// Database 4: Auth/Users (if needed)
		users, _ := r.queryDB("golang",
			`SELECT username FROM users WHERE id = ?`,
			port["id"])

		ports[i]["passengers"] = passengers
		ports[i]["traffic_ticket"] = tickets
		ports[i]["golang"] = users
	}

	return json.Marshal(ports)
}

func convertToPostgresPlaceholders(query string) string {
	counter := 1
	result := ""
	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			result += fmt.Sprintf("$%d", counter)
			counter++
		} else {
			result += string(query[i])
		}
	}
	return result
}

func (r *LautSQLXRepository) queryDB(dbName, query string, args ...interface{}) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	db := r.getDB(dbName)
	if db.DriverName() == "postgres" {
		query = convertToPostgresPlaceholders(query)
	}
	rows, err := db.QueryxContext(ctx, query, args...)
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

		// Convert []byte to string
		for key, value := range row {
			if b, ok := value.([]byte); ok {
				row[key] = string(b)
			}
		}
		results = append(results, row)
	}

	return results, nil
}

func (r *LautSQLXRepository) getDBDriver(dbName string) string {
	db := r.getDB(dbName)
	return db.DriverName()
}

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
