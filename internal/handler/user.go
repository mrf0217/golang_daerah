package handler

import (
	"encoding/json"
	"golang_daerah/internal/service"
	"golang_daerah/pkg/jwtutil"
	"golang_daerah/pkg/response"
	"net/http"
)

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{Service: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteMethodNotAllowed(w)
		return
	}

	var creds service.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		response.WriteBadRequest(w, "Invalid request body: "+err.Error())
		return
	}

	if err := h.Service.Register(creds); err != nil {
		response.WriteBadRequest(w, err.Error())
		return
	}

	response.WriteSuccessResponseCreated(w, []interface{}{}, "Registration successful")
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteMethodNotAllowed(w)
		return
	}

	var creds service.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		response.WriteBadRequest(w, "Invalid request body: "+err.Error())
		return
	}

	token, err := h.Service.Login(creds)
	if err != nil {
		response.WriteUnauthorized(w, err.Error())
		return
	}

	response.WriteSuccessResponseOK(w, map[string]string{"token": token}, "Login successful")
}

func (h *AuthHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		response.WriteUnauthorized(w, "Authorization header required")
		return
	}

	username, err := jwtutil.VerifyToken(authHeader)
	if err != nil {
		response.WriteUnauthorized(w, "Invalid or expired token")
		return
	}

	response.WriteSuccessResponseOK(w, map[string]string{"username": username}, "Welcome, "+username+"!")
}
