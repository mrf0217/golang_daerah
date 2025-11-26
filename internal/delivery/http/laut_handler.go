package http

import (
	"golang_daerah/internal/repository"
	lautgen "golang_daerah/internal/repository/generated/laut"
	"net/http"
)

type LautHandler struct {
	repo *repository.LautRepository
}

func NewLautHandler(repo *repository.LautRepository) *LautHandler {
	return &LautHandler{repo: repo}
}

// Create handles adding new port/terminal data
func (h *LautHandler) Create(w http.ResponseWriter, r *http.Request) {
	CreateHandler(w, r, h.repo, CreateConfig{
		RequirePOST:        true,
		InvalidBodyMessage: "Invalid request body: ",
		InsertErrorMessage: "Failed to insert port/terminal data: %v",
		SuccessMessage:     "%d port/terminal record(s) created successfully",
	})
}

// GetPaginated handles paginated listing of port/terminal data
func (h *LautHandler) GetPaginated(w http.ResponseWriter, r *http.Request) {
	GetPaginatedHandler[lautgen.Laut](w, r, h.repo, PaginatedConfig{
		DefaultPerPage:  10,
		UseGetPaginated: true, // Use GetPaginated method which orders by id ASC (starting from ID 1)
		ErrorMessage:    "Failed to get port/terminal data: %v",
	})
}
