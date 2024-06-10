package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xray-web/web-check-api/checks"
	"github.com/xray-web/web-check-api/config"
	"github.com/xray-web/web-check-api/handlers"
)

type Server struct {
	conf   config.Config
	mux    *http.ServeMux
	checks *checks.Checks
}

func New(conf config.Config) *Server {
	return &Server{
		conf:   conf,
		mux:    http.NewServeMux(),
		checks: checks.NewChecks(),
	}
}

func (s *Server) routes() {
	s.mux.Handle("/", NotFound(nil))

	s.mux.Handle("GET /health", HealthCheck())

	s.mux.Handle("GET /api/block-lists", handlers.HandleBlockLists())
	s.mux.Handle("GET /api/carbon", handlers.HandleCarbon(s.checks.Carbon))
	s.mux.Handle("GET /api/cookies", handlers.HandleCookies())
	s.mux.Handle("GET /api/dns-server", handlers.HandleDNSServer())
	s.mux.Handle("GET /api/dns", handlers.HandleDNS())
	s.mux.Handle("GET /api/dnssec", handlers.HandleDnsSec())
	s.mux.Handle("GET /api/firewall", handlers.HandleFirewall())
	s.mux.Handle("GET /api/get-ip", handlers.HandleGetIP())
	s.mux.Handle("GET /api/headers", handlers.HandleGetHeaders())
	s.mux.Handle("GET /api/hsts", handlers.HandleHsts())
	s.mux.Handle("GET /api/http-security", handlers.HandleHttpSecurity())
	s.mux.Handle("GET /api/legacy-rank", handlers.HandleLegacyRank())
	s.mux.Handle("GET /api/linked-pages", handlers.HandleGetLinks())
	s.mux.Handle("GET /api/ports", handlers.HandleGetPorts())
	s.mux.Handle("GET /api/quality", handlers.HandleGetQuality())
	s.mux.Handle("GET /api/rank", handlers.HandleGetRank(s.checks.Rank))
	s.mux.Handle("GET /api/redirects", handlers.HandleGetRedirects())
	s.mux.Handle("GET /api/social-tags", handlers.HandleGetSocialTags())
	s.mux.Handle("GET /api/tls", handlers.HandleTLS())
	s.mux.Handle("GET /api/trace-route", handlers.HandleTraceRoute())
}

func (s *Server) Run() error {
	s.routes()

	addr := fmt.Sprintf("%s:%s", s.conf.Host, s.conf.Port)
	log.Printf("Server started, listening on: %v\n", addr)
	return http.ListenAndServe(addr, s.CORS(s.mux))
}
