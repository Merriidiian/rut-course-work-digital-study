package schedule

import (
	"context"
	"errors"
	"log"
	"time"
)

type Saga struct {
	repository *Repository
	client     *Client
}

func (saga *Saga) Execute(ctx context.Context, request LessonRequest) SagaResult {
	id, err := saga.repository.NewID(ctx)
	if err != nil {
		return SagaResult{Message: "Cannot create lesson ID"}
	}
	lesson := Lesson{
		ID:        id,
		TeacherID: request.TeacherID,
		SubjectID: request.SubjectID,
		RoomID:    request.RoomID,
		Group:     request.Group,
		StartsAt:  request.StartsAt,
		EndsAt:    request.EndsAt,
		Status:    "Confirmed",
	}
	if err := saga.client.ValidateReferences(ctx, lesson.TeacherID, lesson.SubjectID); err != nil {
		return SagaResult{LessonID: id, Message: "Teacher or subject not found"}
	}
	students, err := saga.client.Students(ctx)
	if err != nil {
		return SagaResult{LessonID: id, Message: "Cannot load students"}
	}
	size := 0
	for _, student := range students {
		if student.Group == lesson.Group {
			size++
		}
	}
	if size == 0 {
		return SagaResult{LessonID: id, Message: "Student group not found"}
	}
	log.Printf("Lesson %s: booking started", id)
	err = saga.client.Book(ctx, lesson, size)
	if err == nil && request.FailAfterBooking {
		err = errors.New("Test failure after booking")
	}
	if err == nil {
		err = saga.repository.Save(ctx, lesson)
	}
	if err != nil {
		rollbackContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		bookingError := saga.client.CancelBooking(rollbackContext, id)
		_, lessonError := saga.repository.Delete(rollbackContext, id)
		compensated := bookingError == nil && lessonError == nil
		log.Printf("Lesson %s: compensation completed = %t", id, compensated)
		return SagaResult{
			LessonID:             id,
			Message:              err.Error(),
			CompensationExecuted: compensated,
		}
	}
	log.Printf("Lesson %s: created", id)
	return SagaResult{Success: true, LessonID: id, Message: "Lesson created"}
}
