package handler

import (
	"golang_daerah/internal/service"
	"golang_daerah/pkg/response"
	"io"
	"net/http"
	"strconv"
)

type LautHandler struct {
	service *service.LautService
}

func NewLautHandler(service *service.LautService) *LautHandler {
	return &LautHandler{service: service}
}

func (h *LautHandler) LautGetCompleteDataHandler(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	data, err := h.service.GetCompleteData(perPage, offset)
	if err != nil {
		response.WriteInternalServerError(w, "Failed to get complete data: "+err.Error())
		return
	}

	// var data []map[string]interface{}
	// if err := json.Unmarshal(jsonData, &data); err != nil {
	// 	response.WriteInternalServerError(w, "Failed to parse response: "+err.Error())
	// 	return
	// }

	response.WritePaginatedResponse(w, data, page, perPage, "Complete data retrieved successfully")
}

func (h *LautHandler) GetPaginated(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        response.WriteBadRequest(w, "Invalid form data")
        return
    }

    page, _ := strconv.Atoi(r.Form.Get("page"))
    perPage, _ := strconv.Atoi(r.Form.Get("perPage"))
    
    if page <= 0 {
        page = 1
    }
    if perPage <= 0 {
        perPage = 10
    }
    
    offset := (page - 1) * perPage

    // Build filters
    filters := make(map[string]string)
    for key, values := range r.Form {
        if key != "page" && key != "perPage" && len(values) > 0 {
            filters[key] = values[0]
        }
    }

    data, err := h.service.GetPaginatedWithFilters(perPage, offset, filters)
    if err != nil {
        response.WriteInternalServerError(w, "Failed to get terminals: "+err.Error())
        return
    }

    response.WritePaginatedResponse(w, data, page, perPage, "Complete data retrieved successfully")
}

	// var data []map[string]interface{}
	// if err := json.Unmarshal(jsonData, &data); err != nil {
	// 	WriteInternalServerError(w, "Failed to parse response: "+err.Error())
	// 	return
	// }



func (h *LautHandler) Create(w http.ResponseWriter, r *http.Request) {
	// dbName := extractDBName(r.URL.Path)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteBadRequest(w, "Invalid request body")
		return
	}

	if err := h.service.Create(body); err != nil {
		response.WriteInternalServerError(w, "Failed to insert terminals: "+err.Error())
		return
	}

	response.WriteSuccessResponseCreated(w, []interface{}{}, "Terminals created successfully")
}
