package buildings

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"university/contracts"
)

type RoomsClient struct {
	http    *http.Client
	baseURL string
}

func NewRoomsClient(baseURL string) *RoomsClient {
	return &RoomsClient{
		http:    &http.Client{Timeout: 5 * time.Second},
		baseURL: baseURL,
	}
}

func (client *RoomsClient) HasRooms(ctx context.Context, buildingID string) (bool, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet,
		client.baseURL+"/api/rooms?buildingId="+url.QueryEscape(buildingID), nil)
	if err != nil {
		return false, err
	}
	response, err := client.http.Do(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Rooms service returned %d", response.StatusCode)
	}
	var rooms []contracts.Room
	if err := json.NewDecoder(response.Body).Decode(&rooms); err != nil {
		return false, err
	}
	return len(rooms) > 0, nil
}
