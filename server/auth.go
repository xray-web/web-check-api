package server

import (
	"net/http"
	"strings"

	"github.com/xray-web/web-check-api/config"
)

type User struct {
	ID    string
	Email string
	Name  string
	Roles []string
}

type Auth struct {
	conf config.Config
	// connection / sdk to auth provider, to trade token for user session token
}

func NewAuth(conf config.Config) *Auth {
	// TODO: reduce scope of conf when we know what auth provider we will use
	return &Auth{conf: conf}
}

func (a *Auth) Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.conf.AuthProvider == "none" {
			h.ServeHTTP(w, r)
			return
		}
		authHeader := r.Header.Get("Authorization")
		// expect "Bearer token" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// use token to get user ID from auth provider
		// TODO: swap token for user session token

	})
}
