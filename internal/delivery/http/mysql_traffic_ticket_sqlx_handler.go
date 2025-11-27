package http

import (
    "io"
    "net/http"
    "strconv"
    "golang_daerah/internal/repository"
)

type MySQLTrafficTicketSQLXHandler struct {
    repo *repository.MySQLTrafficTicketSQLXRepository
}

func NewMySQLTrafficTicketSQLXHandler(repo *repository.MySQLTrafficTicketSQLXRepository) *MySQLTrafficTicketSQLXHandler {
    return &MySQLTrafficTicketSQLXHandler{repo: repo}
}

func (h *MySQLTrafficTicketSQLXHandler) GetPaginated(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))
    if page <= 0 {
        page = 1
    }
    if perPage <= 0 {
        perPage = 10
    }

    offset := (page - 1) * perPage
    jsonData, err := h.repo.GetPaginatedJSON(perPage, offset)
    if err != nil {
        WriteInternalServerError(w, "Failed to get tickets: "+err.Error())
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
}

func (h *MySQLTrafficTicketSQLXHandler) Create(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        WriteBadRequest(w, "Invalid request body")
        return
    }

    if err := h.repo.InsertJSON(body); err != nil {
        WriteInternalServerError(w, "Failed to insert tickets: "+err.Error())
        return
    }

    WriteSuccessResponseCreated(w, []interface{}{}, "Tickets created successfully")
}