package controllers

import (
	"encoding/json"
	"net/http"
)

type ResponseError struct {
	Error string `json:"error"`
}

func JSONError(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ResponseError{
		Error: err,
	})
}

type KV map[string]any

func JSON(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
