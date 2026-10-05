package rooms

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"university/internal/httpjson"
)

type Booking struct {
	ID       string    `json:"id"`
	RoomID   string    `json:"roomId"`
	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`
	Size     int       `json:"size"`
}

func (handler *Handler) RegisterBookings(router *http.ServeMux) {
	router.HandleFunc("POST /api/bookings", handler.book)
	router.HandleFunc("DELETE /api/bookings/{id}", handler.cancelBooking)
}

func (repository *Repository) Book(ctx context.Context, booking Booking) error {
	transaction, err := repository.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	var capacity int
	err = transaction.QueryRowContext(ctx,
		"SELECT capacity FROM rooms WHERE id = $1 FOR UPDATE", booking.RoomID).Scan(&capacity)
	if err != nil {
		return err
	}
	if booking.Size > capacity {
		return errors.New("Room capacity exceeded")
	}
	_, err = transaction.ExecContext(ctx,
		"INSERT INTO bookings (id, room_id, starts_at, ends_at) VALUES ($1, $2, $3, $4)",
		booking.ID, booking.RoomID, booking.StartsAt, booking.EndsAt)
	if err != nil {
		return err
	}
	return transaction.Commit()
}

func (handler *Handler) book(w http.ResponseWriter, request *http.Request) {
	var booking Booking
	if json.NewDecoder(request.Body).Decode(&booking) != nil || booking.ID == "" ||
		booking.Size <= 0 || !booking.EndsAt.After(booking.StartsAt) {
		httpjson.Error(w, 400, "Invalid booking")
		return
	}
	err := handler.repository.Book(request.Context(), booking)
	if errors.Is(err, sql.ErrNoRows) {
		httpjson.Error(w, 404, "Room not found")
		return
	}
	if err != nil {
		httpjson.Error(w, 409, "Room is occupied or capacity exceeded")
		return
	}
	httpjson.Write(w, 201, booking)
}

func (handler *Handler) cancelBooking(w http.ResponseWriter, request *http.Request) {
	_, err := handler.repository.database.ExecContext(request.Context(),
		"DELETE FROM bookings WHERE id = $1", request.PathValue("id"))
	if err != nil {
		httpjson.Error(w, 400, "Invalid booking ID")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
