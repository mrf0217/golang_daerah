package http

import (
	"net/http"

	"golang_daerah/internal/entities"
	"golang_daerah/internal/repository"
)

type TrafficTicketHandler struct {
	repo *repository.TrafficTicketRepository
}

func NewTrafficTicketHandler(repo *repository.TrafficTicketRepository) *TrafficTicketHandler {
	return &TrafficTicketHandler{repo: repo}
}

func (h *TrafficTicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	CreateHandler(w, r, h.repo, CreateConfig{
		RequirePOST:        false,
		InvalidBodyMessage: "Invalid request body: ",
		InsertErrorMessage: "Failed to insert ticket: %v",
		SuccessMessage:     "%d ticket(s) created successfully",
	})
}

// List handles listing tickets with JSON body controlling pagination
// Expected body: {"limit": 10, "offset": 0}
func (h *TrafficTicketHandler) List(w http.ResponseWriter, r *http.Request) {
	ListHandler(w, r, h.repo, ListConfig{
		ReturnPaginated: false,
		ErrorMessage:    "Failed to list tickets: %v",
	})
}

// GetPaginated handles paginated ticket listing with query parameters
func (h *TrafficTicketHandler) GetPaginated(w http.ResponseWriter, r *http.Request) {
	GetPaginatedHandler[entities.TrafficTicket](w, r, h.repo, PaginatedConfig{
		DefaultPerPage:  20,
		UseGetPaginated: true,
		ErrorMessage:    "Failed to get tickets: %v",
	})
}
