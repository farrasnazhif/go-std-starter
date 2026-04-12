package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	// prevent ddos, max 1 MB for the body
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string, err string) error {
	type envelope struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}

	return writeJSON(w, status, &envelope{Status: status, Message: message, Error: err})
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, message string, data any) error {
	type envelope struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}

	return writeJSON(w, status, &envelope{Status: status, Message: message, Data: data})
}
