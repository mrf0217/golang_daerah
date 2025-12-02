package http

// Request Flow Link:
// main.go wires these handlers into the router via routes.Public/Protected, so every incoming HTTP
// request for /api/register or /api/login is processed here before delegating to services/repositories.

import (
	"encoding/json"
	"net/http"

	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang_daerah/config"
	"golang_daerah/internal/entities"

	"golang_daerah/pkg/jwtutil"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

type UserHandler struct {
	Service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (r *UserRepository) CreateUser(username, passwordHash string) error {
	// Create context with timeout to prevent queries from hanging indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	query := `INSERT INTO users (username, password) VALUES ($1, $2) ON CONFLICT (username) DO NOTHING RETURNING id;`
	var id int
	err := r.DB.QueryRowContext(ctx, query, username, passwordHash).Scan(&id)
	if err == sql.ErrNoRows {
		return errors.New("username already exists")
	}
	return handleQueryError(err)
}

func (r *UserRepository) GetUserByUsername(username string) (*entities.User, error) {
	// Create context with timeout to prevent queries from hanging indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	query := `SELECT id, username, password FROM users WHERE username = $1`
	user := &entities.User{}
	err := r.DB.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)

	if err == sql.ErrNoRows {
		fmt.Println("DEBUG: no user found for username:", username)
		return nil, nil // return nil user, no hard error
	} else if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("DEBUG: query timeout for username:", username)
		} else {
			fmt.Println("DEBUG: query error:", err)
		}
		return nil, handleQueryError(err)
	}

	fmt.Println("DEBUG: found user:", user.Username)
	return user, nil
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
