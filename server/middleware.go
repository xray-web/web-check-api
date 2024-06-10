package server

import (
	"encoding/json"
	"net/http"
)

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: determine allowed origins
		w.Header().Set("Access-Control-Allow-Origin", "http://artwork.local:3000")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, api_key, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.Write([]byte{})
			return
		}
		h.ServeHTTP(w, r)
	})
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
