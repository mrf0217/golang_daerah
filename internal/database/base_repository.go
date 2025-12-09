package database

import (
	"context"
	"fmt"
	"golang_daerah/config"

	"github.com/jmoiron/sqlx"
)

// BaseMultiDBRepository provides reusable multi-database functionality
type BaseMultiDBRepository struct {
	dbs map[string]*sqlx.DB
}

// NewBaseMultiDBRepository creates a base repository with multiple database connections
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

// QueryDB executes a query on a specific database
func (r *BaseMultiDBRepository) QueryDB(dbName, query string, args ...interface{}) ([]map[string]interface{}, error) {
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

func (r *BaseMultiDBRepository) InsertDB(dbName, query string, data map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	db := r.getDB(dbName)

	// Handle postgres placeholder conversion if needed
	if db.DriverName() == "postgres" {
		query = convertToPostgresPlaceholders(query)
	}

	_, err := db.NamedExecContext(ctx, query, data)
	if err != nil {
		return handleQueryError(err)
	}
	return nil
}

//example use
// Insert single passenger
// passengerData := map[string]interface{}{
//     "passenger_name": "John Doe",
//     "port_id":        123,
// }
// err := r.insertDB("passenger",
//     `INSERT INTO passenger_plane (passenger_name, port_id)
//      VALUES (:passenger_name, :port_id)`,
//     passengerData)

// // Insert multiple passengers (loop)
// passengers := []map[string]interface{}{
//     {"passenger_name": "John", "port_id": 1},
//     {"passenger_name": "Jane", "port_id": 2},
// }
// for _, p := range passengers {
//     r.insertDB("passenger",
//         `INSERT INTO passenger_plane (passenger_name, port_id)
//          VALUES (:passenger_name, :port_id)`,
//         p)
// }

// ==================== UPDATE HELPER ====================
// updateDB - Helper for UPDATE queries with named parameters
func (r *BaseMultiDBRepository) UpdateDB(dbName, query string, data map[string]interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	db := r.getDB(dbName)

	// Handle postgres placeholder conversion if needed
	if db.DriverName() == "postgres" {
		query = convertToPostgresPlaceholders(query)
	}

	result, err := db.NamedExecContext(ctx, query, data)
	if err != nil {
		return 0, handleQueryError(err)
	}

	// Get number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

//example use
// Update multiple users with same data
// updateData := map[string]interface{}{
//     "status": "active",
//     "age":    25,
// }
// rowsAffected, err := r.updateDB("golang",
//     `UPDATE users SET status = :status WHERE age > :age`,
//     updateData)

// fmt.Printf("Updated %d rows\n", rowsAffected)

// // Update multiple users individually (loop)
// updates := []map[string]interface{}{
//     {"id": 1, "status": "active"},
//     {"id": 2, "status": "inactive"},
// }
// for _, u := range updates {
//     r.updateDB("golang",
//         `UPDATE users SET status = :status WHERE id = :id`,
//         u)
// }

// ==================== DELETE HELPER ====================
// deleteDB - Helper for DELETE queries with positional parameters
func (r *BaseMultiDBRepository) DeleteDB(dbName, query string, args ...interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	db := r.getDB(dbName)

	// Handle postgres placeholder conversion if needed
	if db.DriverName() == "postgres" {
		query = convertToPostgresPlaceholders(query)
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, handleQueryError(err)
	}

	// Get number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

//example use
// Update multiple users with same data
// updateData := map[string]interface{}{
//     "status": "active",
//     "age":    25,
// }
// rowsAffected, err := r.updateDB("golang",
//     `UPDATE users SET status = :status WHERE age > :age`,
//     updateData)

// fmt.Printf("Updated %d rows\n", rowsAffected)

// // Update multiple users individually (loop)
// updates := []map[string]interface{}{
//     {"id": 1, "status": "active"},
//     {"id": 2, "status": "inactive"},
// }
// for _, u := range updates {
//     r.updateDB("golang",
//         `UPDATE users SET status = :status WHERE id = :id`,
//         u)
// }

func (r *BaseMultiDBRepository) getDBDriver(dbName string) string {
	db := r.getDB(dbName)
	return db.DriverName()
}

// Get database by name
func (r *BaseMultiDBRepository) getDB(dbName string) *sqlx.DB {
	if db, exists := r.dbs[dbName]; exists {
		return db
	}
	return r.dbs["default"]
}

func handleQueryError(err error) error {
	if err == context.DeadlineExceeded {
		return fmt.Errorf("database query timeout: request took too long")
	}
	return err
}
