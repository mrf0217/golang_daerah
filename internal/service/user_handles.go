package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"golang_daerah/config"
	"golang_daerah/internal/database"
	"golang_daerah/pkg/jwtutil"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NEW: Multi-DB User Repository
type UserRepository struct {
	*database.BaseMultiDBRepository
}

// NEW: Constructor using multi-DB pattern
func NewUserRepository() *UserRepository {
	// Initialize databases for auth
	dbConfigs := map[string]*sqlx.DB{
		"default": config.InitGolangDBX(),
		"auth":    config.InitAuthDBX(),
	}

	return &UserRepository{
		BaseMultiDBRepository: &database.BaseMultiDBRepository{
			dbs: dbConfigs,
		},
	}
}

// HARDCODED: Configure which databases to use for User operations
// Maintainers can easily modify this function to add/remove databases
// func initializeUserDatabases() map[string]*sqlx.DB {
// 	dbs := make(map[string]*sqlx.DB)

// 	// Default database for users (main auth database)
// 	dbs["default"] = config.InitGolangDBX()

// 	// OPTIONAL: Add backup/replica databases
// 	dbs["auth"] = config.InitAuthDBX()
// 	dbs["mysql"] = config.InitMySQLDBX()

// 	// OPTIONAL: Add other databases if user data needs to be synced
// 	// dbs["passenger"] = config.InitPassengerPlaneDBX()
// 	// dbs["traffic"] = config.InitTrafficDBX()

// 	return dbs
// }

// CreateUser - Now supports multi-DB insert
func (r *UserRepository) CreateUser(username, passwordHash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	// Insert into main database (default)
	db := r.getDB("default")
	query := `INSERT INTO users (username, password) VALUES ($1, $2) ON CONFLICT (username) DO NOTHING RETURNING id;`
	var id int
	err := db.QueryRowContext(ctx, query, username, passwordHash).Scan(&id)
	if err == sql.ErrNoRows {
		return errors.New("username already exists")
	}
	if err != nil {
		return handleQueryError(err)
	}

	// HARDCODED: Automatically replicate to other databases
	// Maintainers can easily add/remove databases here
	// userData := map[string]interface{}{
	// 	"username": username,
	// 	"password": passwordHash,
	// }

	// // Replicate to auth database
	// r.insertDB("auth",
	// 	`INSERT INTO users (username, password) VALUES (:username, :password)`,
	// 	userData)

	// // Replicate to mysql database
	// r.insertDB("mysql",
	// 	`INSERT INTO users (username, password) VALUES (:username, :password)`,
	// 	userData)

	// OPTIONAL: Add more replications here
	// r.insertDB("passenger", `INSERT INTO users ...`, userData)
	// r.insertDB("traffic", `INSERT INTO users ...`, userData)

	return nil
}

// GetUserByUsername - Now supports multi-DB query with fallback
func (r *UserRepository) GetUserByUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	// Try main database first
	db := r.getDB("default")
	query := `SELECT id, username, password FROM users WHERE username = $1`
	user := User{}
	err := db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)

	if err == sql.ErrNoRows {
		// HARDCODED: Fallback to other databases if not found in main
		// Try auth database
		// 	authDB := r.getDB("auth")
		// 	err = authDB.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)
		// 	if err == nil {
		// 		fmt.Println("DEBUG: found user in auth database:", user.Username)
		// 		return &user, nil
		// 	}

		// 	// Try mysql database
		// 	mysqlDB := r.getDB("mysql")
		// 	err = mysqlDB.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)
		// 	if err == nil {
		// 		fmt.Println("DEBUG: found user in mysql database:", user.Username)
		// 		return &user, nil
		// 	}

		// 	fmt.Println("DEBUG: no user found for username:", username)
		// 	return nil, nil
		// } else if err != nil {
		// 	if err == context.DeadlineExceeded {
		// 		fmt.Println("DEBUG: query timeout for username:", username)
		// 	} else {
		// 		fmt.Println("DEBUG: query error:", err)
		// 	}
		return nil, handleQueryError(err)
	}

	fmt.Println("DEBUG: found user in default database:", user.Username)
	return &user, nil
}

// NEW: Get user from ALL databases (for admin/debugging)
// func (r *UserRepository) GetUserFromAllDatabases(username string) (map[string]*User, error) {
// 	results := make(map[string]*User)

// 	// HARDCODED: Query all configured databases
// 	dbNames := []string{"default", "auth", "mysql"}

// 	for _, dbName := range dbNames {
// 		ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
// 		defer cancel()

// 		db := r.getDB(dbName)
// 		query := `SELECT id, username, password FROM users WHERE username = $1`
// 		user := User{}
// 		err := db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)

// 		if err == nil {
// 			results[dbName] = &user
// 		}
// 	}

// 	return results, nil
// }

// NEW: Update user in multiple databases
func (r *UserRepository) UpdateUser(userID int, newPassword string) error {
	updateData := map[string]interface{}{
		"id":       userID,
		"password": newPassword,
	}

	// HARDCODED: Update in all databases
	// Maintainers can easily add/remove databases

	// Update in default database
	_, err := r.updateDB("default",
		`UPDATE users SET password = :password WHERE id = :id`,
		updateData)
	if err != nil {
		return err
	}

	// Update in auth database
	// r.updateDB("auth",
	// 	`UPDATE users SET password = :password WHERE id = :id`,
	// 	updateData)

	// // Update in mysql database
	// r.updateDB("mysql",
	// 	`UPDATE users SET password = :password WHERE id = :id`,
	// 	updateData)

	return nil
}

// NEW: Delete user from multiple databases
func (r *UserRepository) DeleteUser(userID int) error {
	// HARDCODED: Delete from all databases
	// Maintainers can easily add/remove databases

	// Delete from default database
	_, err := r.deleteDB("default",
		`DELETE FROM users WHERE id = ?`,
		userID)
	if err != nil {
		return err
	}

	// Delete from auth database
	// r.deleteDB("auth",
	// 	`DELETE FROM users WHERE id = ?`,
	// 	userID)

	// // Delete from mysql database
	// r.deleteDB("mysql",
	// 	`DELETE FROM users WHERE id = ?`,
	// 	userID)

	return nil
}

type UserHandler struct {
	Service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteMethodNotAllowed(w)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		WriteBadRequest(w, "Invalid request body: "+err.Error())
		return
	}

	if err := h.Service.Register(creds); err != nil {
		WriteBadRequest(w, err.Error())
		return
	}

	WriteSuccessResponseCreated(w, []interface{}{}, "Registration successful (synced to all databases)")
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteMethodNotAllowed(w)
		return
	}

	var creds Credentials
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
