package contracts

import "time"

type Student struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
	Email string `json:"email"`
}

type Teacher struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Subject struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Hours int    `json:"hours"`
}

type Room struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	BuildingID string `json:"buildingId"`
	Capacity   int    `json:"capacity"`
}

type Building struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Lesson struct {
	ID        string    `json:"id"`
	TeacherID string    `json:"teacherId"`
	SubjectID string    `json:"subjectId"`
	RoomID    string    `json:"roomId"`
	Group     string    `json:"group"`
	StartsAt  time.Time `json:"startsAt"`
	EndsAt    time.Time `json:"endsAt"`
	Status    string    `json:"status"`
}

type LessonRequest struct {
	TeacherID        string    `json:"teacherId"`
	SubjectID        string    `json:"subjectId"`
	RoomID           string    `json:"roomId"`
	Group            string    `json:"group"`
	StartsAt         time.Time `json:"startsAt"`
	EndsAt           time.Time `json:"endsAt"`
	FailAfterBooking bool      `json:"failAfterBooking"`
}

type SagaResult struct {
	Success              bool   `json:"success"`
	LessonID             string `json:"lessonId"`
	Message              string `json:"message"`
	CompensationExecuted bool   `json:"compensationExecuted"`
}

type Booking struct {
	ID       string    `json:"id"`
	RoomID   string    `json:"roomId"`
	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`
	Size     int       `json:"size"`
}

type ApiError struct {
	Message string `json:"message"`
}
