package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xray-web/web-check-api/config"
	"github.com/xray-web/web-check-api/handlers"
)

type Server struct {
	conf config.Config
	mux  *http.ServeMux
}

func New(conf config.Config) *Server {
	return &Server{
		conf: conf,
		mux:  http.NewServeMux(),
	}
}

func (s *Server) routes() {
	s.mux.Handle("/", handlers.NotFound(nil))

	s.mux.Handle("GET /headers", handlers.HandleGetHeaders())
	s.mux.Handle("GET /cookies", handlers.HandleCookies())
	s.mux.Handle("GET /carbon", handlers.HandleCarbon())
	s.mux.Handle("GET /block-lists", handlers.HandleBlockLists())
	s.mux.Handle("GET /dns-server", handlers.HandleDNSServer())
	s.mux.Handle("GET /dns", handlers.HandleDNS())
	s.mux.Handle("GET /dnssec", handlers.HandleDnsSec())
	s.mux.Handle("GET /firewall", handlers.HandleFirewall())
	s.mux.Handle("GET /get-ip", handlers.HandleGetIP())
	s.mux.Handle("GET /hsts", handlers.HandleHsts())
	s.mux.Handle("GET /http-security", handlers.HandleHttpSecurity())
	s.mux.Handle("GET /legacy-rank", handlers.HandleLegacyRank())
	s.mux.Handle("GET /linked-pages", handlers.HandleGetLinks())
	s.mux.Handle("GET /ports", handlers.HandleGetPorts())
	s.mux.Handle("GET /quality", handlers.HandleGetQuality())
	s.mux.Handle("GET /rank", handlers.HandleGetRank())
	s.mux.Handle("GET /redirects", handlers.HandleGetRedirects())
	s.mux.Handle("GET /social-tags", handlers.HandleGetSocialTags())
	s.mux.Handle("GET /tls", handlers.HandleTLS())
	s.mux.Handle("GET /trace-route", handlers.HandleTraceRoute())
}

func (s *Server) Run() error {
	s.routes()

	addr := fmt.Sprintf("%s:%s", s.conf.Host, s.conf.Port)
	log.Printf("Server started, listening on: %v\n", addr)
	return http.ListenAndServe(addr, s.mux)
}
