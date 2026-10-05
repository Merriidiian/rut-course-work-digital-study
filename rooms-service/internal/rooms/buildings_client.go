package rooms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"university/contracts"
)

var errBuildingNotFound = errors.New("Building not found")

type BuildingsClient struct {
	http    *http.Client
	baseURL string
}

func NewBuildingsClient(baseURL string) *BuildingsClient {
	return &BuildingsClient{
		http:    &http.Client{Timeout: 5 * time.Second},
		baseURL: baseURL,
	}
}

func (client *BuildingsClient) Find(ctx context.Context, id string) (contracts.Building, error) {
	var building contracts.Building
	request, err := http.NewRequestWithContext(ctx, http.MethodGet,
		client.baseURL+"/api/buildings/"+id, nil)
	if err != nil {
		return building, err
	}
	response, err := client.http.Do(request)
	if err != nil {
		return building, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound {
		return building, errBuildingNotFound
	}
	if response.StatusCode != http.StatusOK {
		return building, fmt.Errorf("Buildings service returned %d", response.StatusCode)
	}
	err = json.NewDecoder(response.Body).Decode(&building)
	return building, err
}
