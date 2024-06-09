package handlers

import (
	"net/http"
)

// HandleHealthCheck returns the status of the application
func HandleHealthCheck() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok", "message":"We're alive!"}`))
	})
}
