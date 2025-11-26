package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang_daerah/config"
	"golang_daerah/internal/entities"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
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
