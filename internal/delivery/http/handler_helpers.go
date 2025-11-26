package http

// Request Flow Link:
// main.go registers handlers (routes.Public/Protected) that live in internal/delivery/http.
// Those concrete handlers delegate shared logic to the generic helpers in this file, meaning every
// HTTP request routed from main.go eventually passes through these helpers before reaching repositories.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Repository interfaces for generic handlers

// InsertableRepository defines the interface for repositories that can insert entities
type InsertableRepository[T any] interface {
	Insert(item *T) error
}

// ListableRepository defines the interface for repositories that can list entities
type ListableRepository[T any] interface {
	List(limit, offset int) ([]*T, error)
}

// PaginatableRepository defines the interface for repositories that can paginate entities
type PaginatableRepository[T any] interface {
	GetPaginated(limit, offset int) ([]T, error)
	List(limit, offset int) ([]*T, error)
}

// PaginatableRepositoryWithGetPaginated is for repos that have GetPaginated method
type PaginatableRepositoryWithGetPaginated[T any] interface {
	GetPaginated(limit, offset int) ([]T, error)
}

// CreateConfig provides configuration for CreateHandler
type CreateConfig struct {
	RequirePOST        bool
	InvalidBodyMessage string
	InsertErrorMessage string
	SuccessMessage     string
}

// CreateHandler handles creation of entities (array or single)
func CreateHandler[T any](
	w http.ResponseWriter,
	r *http.Request,
	repo InsertableRepository[T],
	config CreateConfig,
) {
	if config.RequirePOST && r.Method != http.MethodPost {
		WriteMethodNotAllowed(w)
		return
	}

	var items []T
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		var single T
		if err2 := json.NewDecoder(r.Body).Decode(&single); err2 != nil {
			WriteBadRequest(w, config.InvalidBodyMessage+err.Error())
			return
		}
		items = append(items, single)
	}

	for _, item := range items {
		if err := repo.Insert(&item); err != nil {
			WriteInternalServerError(w, fmt.Sprintf(config.InsertErrorMessage, err))
			return
		}
	}

	WriteSuccessResponseCreated(w, []interface{}{}, fmt.Sprintf(config.SuccessMessage, len(items)))
}

// PaginatedConfig provides configuration for GetPaginatedHandler
type PaginatedConfig struct {
	DefaultPerPage  int
	UseGetPaginated bool
	ErrorMessage    string
}

// GetPaginatedHandler handles paginated listing with query parameters or headers
// Supports both ?page=1&perPage=20 in URL or X-Page and X-Per-Page headers
func GetPaginatedHandler[T any](
	w http.ResponseWriter,
	r *http.Request,
	repo interface{},
	config PaginatedConfig,
) {
	// Try query parameters first, then fall back to headers
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = r.Header.Get("Page")
	}

	perPageStr := r.URL.Query().Get("perPage")
	if perPageStr == "" {
		perPageStr = r.Header.Get("PerPage")
	}

	page, _ := strconv.Atoi(pageStr)
	perPage, _ := strconv.Atoi(perPageStr)
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = config.DefaultPerPage
	}

	offset := (page - 1) * perPage

	var items interface{}
	var err error

	if config.UseGetPaginated {
		// Type assert to get GetPaginated method
		if paginatedRepo, ok := repo.(PaginatableRepositoryWithGetPaginated[T]); ok {
			items, err = paginatedRepo.GetPaginated(perPage, offset)
		} else {
			// Fallback to List if GetPaginated not available
			if listRepo, ok := repo.(ListableRepository[T]); ok {
				items, err = listRepo.List(perPage, offset)
			} else {
				WriteInternalServerError(w, "Repository does not support pagination")
				return
			}
		}
	} else {
		// Use List method
		if listRepo, ok := repo.(ListableRepository[T]); ok {
			items, err = listRepo.List(perPage, offset)
		} else {
			WriteInternalServerError(w, "Repository does not support listing")
			return
		}
	}

	if err != nil {
		WriteInternalServerError(w, fmt.Sprintf(config.ErrorMessage, err))
		return
	}

	WritePaginatedResponse(w, items, page, perPage, "OK")
}

// ListConfig provides configuration for ListHandler
type ListConfig struct {
	ReturnPaginated bool
	ErrorMessage    string
}

// ListHandler handles listing with JSON body pagination
func ListHandler[T any](
	w http.ResponseWriter,
	r *http.Request,
	repo ListableRepository[T],
	config ListConfig,
) {
	if r.Method != http.MethodPost {
		WriteMethodNotAllowed(w)
		return
	}

	var payload struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteBadRequest(w, "Invalid request body: "+err.Error())
		return
	}

	items, err := repo.List(payload.Limit, payload.Offset)
	if err != nil {
		WriteInternalServerError(w, fmt.Sprintf(config.ErrorMessage, err))
		return
	}

	if config.ReturnPaginated {
		page := payload.Offset/payload.Limit + 1
		if payload.Limit == 0 {
			page = 1
		}
		WritePaginatedResponse(w, items, page, payload.Limit, "")
	} else {
		WriteSuccessResponseOK(w, items, "")
	}
}
