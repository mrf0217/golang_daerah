package http

import (
	"net/http"

	"golang_daerah/internal/entities"
	"golang_daerah/internal/repository"
)

type SQLServerTrafficTicketHandler struct {
	repo *repository.SQLServerTrafficTicketRepository
}

func NewSQLServerTrafficTicketHandler(repo *repository.SQLServerTrafficTicketRepository) *SQLServerTrafficTicketHandler {
	return &SQLServerTrafficTicketHandler{repo: repo}
}

func (h *SQLServerTrafficTicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	CreateHandler(w, r, h.repo, CreateConfig{
		RequirePOST:        true,
		InvalidBodyMessage: "Invalid request body: ",
		InsertErrorMessage: "Failed to insert ticket to SQL Server: %v",
		SuccessMessage:     "%d ticket(s) created successfully in SQL Server",
	})
}

func (h *SQLServerTrafficTicketHandler) List(w http.ResponseWriter, r *http.Request) {
	ListHandler(w, r, h.repo, ListConfig{
		ReturnPaginated: false,
		ErrorMessage:    "Failed to list tickets from SQL Server: %v",
	})
}

func (h *SQLServerTrafficTicketHandler) GetPaginated(w http.ResponseWriter, r *http.Request) {
	GetPaginatedHandler[entities.TrafficTicket](w, r, h.repo, PaginatedConfig{
		DefaultPerPage:  10,
		UseGetPaginated: false,
		ErrorMessage:    "Failed to get tickets from SQL Server: %v",
	})
}
