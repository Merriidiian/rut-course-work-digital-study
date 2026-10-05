package rooms

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"university/internal/httpjson"
)

type Handler struct {
	repository *Repository
}

func NewHandler(database *sql.DB) *Handler {
	return &Handler{repository: &Repository{database: database}}
}

func (handler *Handler) Register(router *http.ServeMux) {
	router.HandleFunc("GET /api/rooms", handler.list)
	router.HandleFunc("GET /api/rooms/{id}", handler.find)
	router.HandleFunc("POST /api/rooms", handler.save)
	router.HandleFunc("PUT /api/rooms/{id}", handler.save)
	router.HandleFunc("DELETE /api/rooms/{id}", handler.delete)
	handler.RegisterBookings(router)
}

func (handler *Handler) list(w http.ResponseWriter, request *http.Request) {
	entities, err := handler.repository.List(request.Context())
	if err != nil {
		httpjson.Error(w, 500, "Cannot load rooms")
		return
	}
	httpjson.Write(w, 200, entities)
}

func (handler *Handler) find(w http.ResponseWriter, request *http.Request) {
	entity, err := handler.repository.Find(request.Context(), request.PathValue("id"))
	if errors.Is(err, sql.ErrNoRows) {
		httpjson.Error(w, 404, "Not found")
		return
	}
	if err != nil {
		httpjson.Error(w, 400, "Invalid ID")
		return
	}
	httpjson.Write(w, 200, entity)
}

func (handler *Handler) save(w http.ResponseWriter, request *http.Request) {
	var entity Room
	if json.NewDecoder(request.Body).Decode(&entity) != nil || strings.TrimSpace(entity.Name) == "" || entity.BuildingID == "" || entity.Capacity <= 0 {
		httpjson.Error(w, 400, "Invalid rooms")
		return
	}
	entity.ID = request.PathValue("id")
	saved, err := handler.repository.Save(request.Context(), entity)
	if errors.Is(err, sql.ErrNoRows) {
		httpjson.Error(w, 404, "Not found")
		return
	}
	if err != nil {
		httpjson.Error(w, 409, "Cannot save rooms")
		return
	}
	status := http.StatusOK
	if request.Method == http.MethodPost {
		status = http.StatusCreated
	}
	httpjson.Write(w, status, saved)
}

func (handler *Handler) delete(w http.ResponseWriter, request *http.Request) {
	deleted, err := handler.repository.Delete(request.Context(), request.PathValue("id"))
	if err != nil {
		httpjson.Error(w, 409, "Record is referenced or ID is invalid")
		return
	}
	if !deleted {
		httpjson.Error(w, 404, "Not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
