package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

var ErrMissingURLParameter = errors.New("missing URL parameter")
var ErrInvalidURL = errors.New("invalid URL")

type ResponseError struct {
	Error string `json:"error"`
}

func JSONError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ResponseError{
		Error: err.Error(),
	})
}

type KV map[string]any

func JSON(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func NotFound(h http.Handler) http.Handler {
	type Response struct {
		Status string `json:"status"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h == nil || r.URL.Path != "/" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Response{Status: "route not found"})
			return
		}
		h.ServeHTTP(w, r)
	})
}
