package schedule

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	http     *http.Client
	coreURL  string
	roomsURL string
}

func NewClient(coreURL, roomsURL string) *Client {
	return &Client{
		http:     &http.Client{Timeout: 5 * time.Second},
		coreURL:  coreURL,
		roomsURL: roomsURL,
	}
}

func (client *Client) request(ctx context.Context, method, url string, body, result any) error {
	var data bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&data).Encode(body); err != nil {
			return err
		}
	}
	request, err := http.NewRequestWithContext(ctx, method, url, &data)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.http.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode >= 300 {
		return fmt.Errorf("Service returned status %d", response.StatusCode)
	}
	if result != nil {
		return json.NewDecoder(response.Body).Decode(result)
	}
	return nil
}

func (client *Client) Students(ctx context.Context) ([]Student, error) {
	students := []Student{}
	err := client.request(ctx, http.MethodGet, client.coreURL+"/api/students", nil, &students)
	return students, err
}

func (client *Client) Student(ctx context.Context, id string) (Student, error) {
	var student Student
	err := client.request(ctx, http.MethodGet, client.coreURL+"/api/students/"+id, nil, &student)
	return student, err
}

func (client *Client) ValidateReferences(ctx context.Context, teacherID, subjectID string) error {
	if err := client.request(ctx, http.MethodGet, client.coreURL+"/api/teachers/"+teacherID, nil, nil); err != nil {
		return err
	}
	return client.request(ctx, http.MethodGet, client.coreURL+"/api/subjects/"+subjectID, nil, nil)
}

func (client *Client) Book(ctx context.Context, lesson Lesson, size int) error {
	booking := struct {
		ID       string    `json:"id"`
		RoomID   string    `json:"roomId"`
		StartsAt time.Time `json:"startsAt"`
		EndsAt   time.Time `json:"endsAt"`
		Size     int       `json:"size"`
	}{lesson.ID, lesson.RoomID, lesson.StartsAt, lesson.EndsAt, size}
	return client.request(ctx, http.MethodPost, client.roomsURL+"/api/bookings", booking, nil)
}

func (client *Client) CancelBooking(ctx context.Context, id string) error {
	return client.request(ctx, http.MethodDelete, client.roomsURL+"/api/bookings/"+id, nil, nil)
}
