package http

// Request Flow Link:
// When main.go routes an HTTP request to a handler, that handler builds its JSON output through
// the helpers in this file so that every response in the main.go flow shares the same structure.

import (
	"encoding/json"
	"net/http"
)

// Response represents the standard API response structure
type Response struct {
	Status  bool        `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
	Page    int         `json:"page,omitempty"`
	PerPage int         `json:"perPage,omitempty"`
}

// WriteErrorResponse writes an error response with the given status code and message
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Status:  false,
		Data:    []interface{}{},
		Message: message,
	})
}

// WriteSuccessResponse writes a success response with the given status code, data, and message
func WriteSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Status:  true,
		Data:    data,
		Message: message,
	})
}

// WriteSuccessResponseOK writes a 200 OK success response (most common)
func WriteSuccessResponseOK(w http.ResponseWriter, data interface{}, message string) {
	WriteSuccessResponse(w, http.StatusOK, data, message)
}

// WriteSuccessResponseCreated writes a 201 Created success response
func WriteSuccessResponseCreated(w http.ResponseWriter, data interface{}, message string) {
	WriteSuccessResponse(w, http.StatusCreated, data, message)
}

// WritePaginatedResponse writes a success response with pagination info
func WritePaginatedResponse(w http.ResponseWriter, data interface{}, page, perPage int, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status:  true,
		Data:    data,
		Page:    page,
		PerPage: perPage,
		Message: message,
	})
}

// Convenience functions for common error cases
func WriteBadRequest(w http.ResponseWriter, message string) {
	WriteErrorResponse(w, http.StatusBadRequest, message)
}

func WriteUnauthorized(w http.ResponseWriter, message string) {
	WriteErrorResponse(w, http.StatusUnauthorized, message)
}

func WriteMethodNotAllowed(w http.ResponseWriter) {
	WriteErrorResponse(w, http.StatusMethodNotAllowed, "Only POST method allowed")
}

func WriteInternalServerError(w http.ResponseWriter, message string) {
	WriteErrorResponse(w, http.StatusInternalServerError, message)
}

