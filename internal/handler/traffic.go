package handler

import (
	"golang_daerah/internal/service"
	"golang_daerah/pkg/response"
	"io"
	"net/http"
	"strconv"
)

type TrafficHandler struct {
	service *service.TrafficService
}

func NewTrafficHandler(service *service.TrafficService) *TrafficHandler {
	return &TrafficHandler{service: service}
}

func (h *TrafficHandler) GetPaginated(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	data, err := h.service.GetPaginated(perPage, offset)
	if err != nil {
		response.WriteInternalServerError(w, "Failed to get tickets: "+err.Error())
		return
	}

	// var data []map[string]interface{}
	// if err := json.Unmarshal(jsonData, &data); err != nil {
	// 	WriteInternalServerError(w, "Failed to parse response: "+err.Error())
	// 	return
	// }

	response.WritePaginatedResponse(w, data, page, perPage, "Complete data retrieved successfully")
}

func (h *TrafficHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteBadRequest(w, "Invalid request body")
		return
	}

	if err := h.service.Create(body); err != nil {
		response.WriteInternalServerError(w, "Failed to insert tickets: "+err.Error())
		return
	}

	response.WriteSuccessResponseCreated(w, []interface{}{}, "Tickets created successfully")
}
