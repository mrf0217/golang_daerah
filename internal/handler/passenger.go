package handler

import (
	"golang_daerah/internal/service"
	"golang_daerah/pkg/response"
	"io"
	"net/http"
	"strconv"
)

type PassengerHandler struct {
	service *service.PassengerPlaneService
}

func NewPassengerHandler(service *service.PassengerPlaneService) *PassengerHandler {
	return &PassengerHandler{service: service}
}

func (h *PassengerHandler) GetPaginated(w http.ResponseWriter, r *http.Request) {
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
		response.WriteInternalServerError(w, "Failed to get passengers: "+err.Error())
		return
	}

	// var data []map[string]interface{}
	// if err := json.Unmarshal(jsonData, &data); err != nil {
	// 	response.WriteInternalServerError(w, "Failed to parse response: "+err.Error())
	// 	return
	// }

	response.WritePaginatedResponse(w, data, page, perPage, "Complete data retrieved successfully")
}

func (h *PassengerHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteBadRequest(w, "Invalid request body")
		return
	}

	if err := h.service.Create(body); err != nil {
		response.WriteInternalServerError(w, "Failed to insert passengers: "+err.Error())
		return
	}

	response.WriteSuccessResponseCreated(w, []interface{}{}, "Passengers created successfully")
}
