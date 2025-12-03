package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type ResponsePaginate struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`

	Meta Meta `json:"meta"`
}

type Meta struct {
	Page        uint `json:"page"`
	TotalPage   uint `json:"total_page"`
	TotalData   uint `json:"total_data"`
	DataperPage uint `json:"per_page"`
}

func ParseJson(r *http.Request, model interface{}) error {
	return json.NewDecoder(r.Body).Decode(model)
}

func ResponseJson(w http.ResponseWriter, status int, response any) {
	jsonBytes, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonBytes)
}
