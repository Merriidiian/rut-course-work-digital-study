package httpjson

import (
	"encoding/json"
	"net/http"

	"university/contracts"
)

func Write(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		return
	}
}

func Error(w http.ResponseWriter, status int, message string) {
	Write(w, status, contracts.ApiError{Message: message})
}
