package buildings

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"university/buildings/internal/httpjson"
)

type Handler struct {
	repository *Repository
	rooms      *RoomsClient
}

func NewHandler(database *sql.DB) *Handler {
	return &Handler{
		repository: &Repository{database: database},
		rooms:      NewRoomsClient(os.Getenv("ROOMS_URL")),
	}
}

func (handler *Handler) Register(router *http.ServeMux) {
	router.HandleFunc("GET /api/buildings", handler.list)
	router.HandleFunc("GET /api/buildings/{id}", handler.find)
	router.HandleFunc("POST /api/buildings", handler.save)
	router.HandleFunc("PUT /api/buildings/{id}", handler.save)
	router.HandleFunc("DELETE /api/buildings/{id}", handler.delete)

}

func (handler *Handler) list(w http.ResponseWriter, request *http.Request) {
	entities, err := handler.repository.List(request.Context())
	if err != nil {
		httpjson.Error(w, 500, "Cannot load buildings")
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
	var entity Building
	if json.NewDecoder(request.Body).Decode(&entity) != nil || strings.TrimSpace(entity.Name) == "" || strings.TrimSpace(entity.Address) == "" {
		httpjson.Error(w, 400, "Invalid buildings")
		return
	}
	entity.ID = request.PathValue("id")
	saved, err := handler.repository.Save(request.Context(), entity)
	if errors.Is(err, sql.ErrNoRows) {
		httpjson.Error(w, 404, "Not found")
		return
	}
	if err != nil {
		httpjson.Error(w, 409, "Cannot save buildings")
		return
	}
	status := http.StatusOK
	if request.Method == http.MethodPost {
		status = http.StatusCreated
	}
	httpjson.Write(w, status, saved)
}

func (handler *Handler) delete(w http.ResponseWriter, request *http.Request) {
	used, err := handler.rooms.HasRooms(request.Context(), request.PathValue("id"))
	if err != nil {
		httpjson.Error(w, 502, "Rooms service is unavailable")
		return
	}
	if used {
		httpjson.Error(w, 409, "Building has rooms")
		return
	}
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
