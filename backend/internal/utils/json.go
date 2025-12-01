package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func ParseJson(r *http.Request, model interface{}) error {
	return json.NewDecoder(r.Body).Decode(model)
}

func ResponseJson(w http.ResponseWriter, status int, response Response) {
	jsonBytes, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonBytes)
}
