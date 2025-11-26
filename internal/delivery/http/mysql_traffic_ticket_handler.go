package http

import (
	"net/http"

	"golang_daerah/internal/entities"
	"golang_daerah/internal/repository"
)

type MySQLTrafficTicketHandler struct {
	repo *repository.MySQLTrafficTicketRepository
}

func NewMySQLTrafficTicketHandler(repo *repository.MySQLTrafficTicketRepository) *MySQLTrafficTicketHandler {
	return &MySQLTrafficTicketHandler{repo: repo}
}

func (h *MySQLTrafficTicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	CreateHandler(w, r, h.repo, CreateConfig{
		RequirePOST:        true,
		InvalidBodyMessage: "Invalid request body: ",
		InsertErrorMessage: "Failed to insert ticket to MySQL: %v",
		SuccessMessage:     "%d ticket(s) created successfully in MySQL",
	})
}

func (h *MySQLTrafficTicketHandler) List(w http.ResponseWriter, r *http.Request) {
	ListHandler(w, r, h.repo, ListConfig{
		ReturnPaginated: true,
		ErrorMessage:    "Failed to list tickets from MySQL: %v",
	})
}

func (h *MySQLTrafficTicketHandler) GetPaginated(w http.ResponseWriter, r *http.Request) {
	GetPaginatedHandler[entities.TrafficTicket](w, r, h.repo, PaginatedConfig{
		DefaultPerPage:  10,
		UseGetPaginated: true, // Use GetPaginated method which orders by id ASC (starting from ID 1)
		ErrorMessage:    "Failed to get tickets from MySQL: %v",
	})
}
