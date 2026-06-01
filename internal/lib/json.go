package lib

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func WriteJSONError(w http.ResponseWriter, status int, message, errMsg string) error {
	type envelope struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	return WriteJSON(w, status, &envelope{Status: status, Message: message, Error: errMsg})
}

func JSONResponse(w http.ResponseWriter, status int, message string, data any) error {
	type envelope struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}
	return WriteJSON(w, status, &envelope{Status: status, Message: message, Data: data})
}
