package http

import (
	"net/http"

	"golang_daerah/internal/entities"
	"golang_daerah/internal/repository"
)

type PassengerPlaneHandler struct {
	repo *repository.PassengerPlaneRepository
}

func NewPassengerPlaneHandler(repo *repository.PassengerPlaneRepository) *PassengerPlaneHandler {
	return &PassengerPlaneHandler{repo: repo}
}

func (h *PassengerPlaneHandler) Create(w http.ResponseWriter, r *http.Request) {
	CreateHandler(w, r, h.repo, CreateConfig{
		RequirePOST:        true,
		InvalidBodyMessage: "Invalid request body: ",
		InsertErrorMessage: "Failed to insert ticket to passenger_plane: %v",
		SuccessMessage:     "%d passenger record(s) created successfully in passenger_plane",
	})
}

func (h *PassengerPlaneHandler) List(w http.ResponseWriter, r *http.Request) {
	ListHandler(w, r, h.repo, ListConfig{
		ReturnPaginated: true,
		ErrorMessage:    "Failed to list tickets from passenger_plane: %v",
	})
}

func (h *PassengerPlaneHandler) GetPaginated(w http.ResponseWriter, r *http.Request) {
	GetPaginatedHandler[entities.Passenger](w, r, h.repo, PaginatedConfig{
		DefaultPerPage:  10,
		UseGetPaginated: true, // Use GetPaginated method which orders by id ASC (starting from ID 1)
		ErrorMessage:    "Failed to get tickets from passenger_plane: %v",
	})
}
