package schedule

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"university/internal/httpjson"
)

type Handler struct {
	repository *Repository
	client     *Client
	saga       *Saga
}

func NewHandler(database *sql.DB) *Handler {
	repository := &Repository{database: database}
	client := NewClient(os.Getenv("CORE_URL"), os.Getenv("ROOMS_URL"))
	return &Handler{
		repository: repository,
		client:     client,
		saga:       &Saga{repository: repository, client: client},
	}
}

func (handler *Handler) Register(router *http.ServeMux) {
	router.HandleFunc("GET /api/lessons", handler.list)
	router.HandleFunc("POST /api/lessons", handler.create)
	router.HandleFunc("DELETE /api/lessons/{id}", handler.delete)
	router.HandleFunc("GET /api/students/{id}/schedule", handler.studentSchedule)
}

func (handler *Handler) list(w http.ResponseWriter, request *http.Request) {
	lessons, err := handler.repository.List(request.Context(),
		request.URL.Query().Get("group"), request.URL.Query().Get("teacherId"))
	if err != nil {
		httpjson.Error(w, 500, "Cannot load lessons")
		return
	}
	httpjson.Write(w, 200, lessons)
}

func (handler *Handler) create(w http.ResponseWriter, request *http.Request) {
	var lesson LessonRequest
	if json.NewDecoder(request.Body).Decode(&lesson) != nil || lesson.TeacherID == "" ||
		lesson.SubjectID == "" || lesson.RoomID == "" || strings.TrimSpace(lesson.Group) == "" ||
		!lesson.EndsAt.After(lesson.StartsAt) {
		httpjson.Error(w, 400, "Invalid lesson")
		return
	}
	result := handler.saga.Execute(request.Context(), lesson)
	status := http.StatusCreated
	if !result.Success {
		status = http.StatusBadRequest
	}
	httpjson.Write(w, status, result)
}

func (handler *Handler) delete(w http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")
	if err := handler.client.CancelBooking(request.Context(), id); err != nil {
		httpjson.Error(w, 502, "Cannot cancel booking")
		return
	}
	deleted, err := handler.repository.Delete(request.Context(), id)
	if err != nil {
		httpjson.Error(w, 400, "Invalid lesson ID")
		return
	}
	if !deleted {
		httpjson.Error(w, 404, "Not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (handler *Handler) studentSchedule(w http.ResponseWriter, request *http.Request) {
	student, err := handler.client.Student(request.Context(), request.PathValue("id"))
	if err != nil {
		httpjson.Error(w, 404, "Student not found")
		return
	}
	lessons, err := handler.repository.List(request.Context(), student.Group, "")
	if err != nil {
		httpjson.Error(w, 500, "Cannot load lessons")
		return
	}
	httpjson.Write(w, 200, lessons)
}
