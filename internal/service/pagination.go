package service

// Request Flow Link:
// Database-facing methods that are invoked when main.go routes reach repository instances
// call into these helpers to normalize pagination parameters and unify timeout error handling.

import (
	"context"
	// "database/sql"
	"errors"
)

// normalizePagination validates and normalizes pagination parameters
// Returns default values if invalid: limit defaults to 10, offset defaults to 0
// func normalizePagination(limit, offset int) (int, int) {
// 	if limit <= 0 {
// 		limit = 10
// 	}
// 	if offset < 0 {
// 		offset = 0
// 	}
// 	return limit, offset
// }

// handleQueryError checks for context timeout and returns appropriate error
func handleQueryError(err error) error {
	if err == context.DeadlineExceeded {
		return errors.New("database query timeout: request took too long")
	}
	return err
}

// scanRows iterates through rows and scans them using the provided scan function
// Returns the scanned items and any error encountered
// func scanRows[T any](rows *sql.Rows, scanFunc func(*sql.Rows) (T, error)) ([]T, error) {
// 	defer rows.Close()

// 	var items []T
// 	for rows.Next() {
// 		item, err := scanFunc(rows)
// 		if err != nil {
// 			return nil, err
// 		}
// 		items = append(items, item)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return items, nil
// }
