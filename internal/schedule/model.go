package schedule

import "time"

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

type Student struct {
	ID    string `json:"id"`
	Group string `json:"group"`
}

type SagaResult struct {
	Success              bool   `json:"success"`
	LessonID             string `json:"lessonId"`
	Message              string `json:"message"`
	CompensationExecuted bool   `json:"compensationExecuted"`
}
