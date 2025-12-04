package http

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
func (r *BaseMultiDBRepository) queryDB(dbName, query string, args ...interface{}) ([]map[string]interface{}, error) {
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






