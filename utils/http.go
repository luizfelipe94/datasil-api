package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ResponseOk(w http.ResponseWriter, statusCode int, data any) error {
	if statusCode >= 400 {
		return errors.New("Invalid status code")
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func ResponseError(w http.ResponseWriter, statusCode int, message string) error {
	if statusCode < 400 {
		return errors.New("Invalid status code")
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func ParseBody(r *http.Request, data any) error {
	if r.Body == nil {
		return errors.New("Invalid body")
	}
	return json.NewDecoder(r.Body).Decode(data)
}
