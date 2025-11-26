package http

// Request Flow Link:
// main.go wires these handlers into the router via routes.Public/Protected, so every incoming HTTP
// request for /api/register or /api/login is processed here before delegating to services/repositories.

import (
	"encoding/json"
	"net/http"

	"golang_daerah/internal/entities"
	"golang_daerah/internal/usecases"
	"golang_daerah/pkg/jwtutil"
)

type UserHandler struct {
	Service *usecases.UserService
}

func NewUserHandler(service *usecases.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteMethodNotAllowed(w)
		return
	}

	var creds entities.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		WriteBadRequest(w, "Invalid request body: "+err.Error())
		return
	}

	if err := h.Service.Register(creds); err != nil {
		WriteBadRequest(w, err.Error())
		return
	}

	WriteSuccessResponseCreated(w, []interface{}{}, "Registration successful")
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteMethodNotAllowed(w)
		return
	}

	var creds entities.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		WriteBadRequest(w, "Invalid request body: "+err.Error())
		return
	}

	token, err := h.Service.Login(creds)
	if err != nil {
		WriteUnauthorized(w, err.Error())
		return
	}

	WriteSuccessResponseOK(w, map[string]string{"token": token}, "Login successful")
}

func (h *UserHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		WriteUnauthorized(w, "Authorization header required")
		return
	}

	username, err := jwtutil.VerifyToken(authHeader)
	if err != nil {
		WriteUnauthorized(w, "Invalid or expired token")
		return
	}

	WriteSuccessResponseOK(w, map[string]string{"username": username}, "Welcome, "+username+"!")
}
