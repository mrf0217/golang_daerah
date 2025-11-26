package usecases

// Request Flow Link:
// In main.go the user handler is constructed with a UserService, so every /api/register or /api/login
// request passes through these methods before hitting repositories or JWT utilities.

import (
	"errors"
	"fmt"
	"golang_daerah/internal/entities"
	"golang_daerah/internal/repository"
	"golang_daerah/pkg/jwtutil"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) Register(creds entities.Credentials) error {
	if creds.Username == "" || creds.Password == "" {
		return errors.New("username and password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.Repo.CreateUser(creds.Username, string(hashedPassword))
}

func (s *UserService) Login(creds entities.Credentials) (string, error) {
	user, err := s.Repo.GetUserByUsername(creds.Username)
	if err != nil {
		fmt.Println("DEBUG: repo error:", err)
		return "", errors.New("internal server error")
	}
	if user == nil {
		fmt.Println("DEBUG: username not found:", creds.Username)
		return "", errors.New("invalid username or password")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)) != nil {
		fmt.Println("DEBUG MISMATCH:", user.PasswordHash, creds.Password)
		return "", errors.New("invalid username or password")
	}

	token, err := jwtutil.GenerateToken(user.Username, time.Hour)
	if err != nil {
		return "", err
	}

	return token, nil
}
