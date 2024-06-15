package server

import (
	"context"
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
	srv    *http.Server
}

func New(conf config.Config) *Server {
	return &Server{
		srv:    &http.Server{},
		conf:   conf,
		mux:    http.NewServeMux(),
		checks: checks.NewChecks(),
	}
}

func (s *Server) routes() {
	s.mux.Handle("/", NotFound(nil))

	s.mux.Handle("GET /health", HealthCheck())

	s.mux.Handle("GET /api/block-lists", handlers.HandleBlockLists(s.checks.BlockList))
	s.mux.Handle("GET /api/carbon", handlers.HandleCarbon(s.checks.Carbon))
	s.mux.Handle("GET /api/cookies", handlers.HandleCookies())
	s.mux.Handle("GET /api/dns-server", handlers.HandleDNSServer())
	s.mux.Handle("GET /api/dns", handlers.HandleDNS())
	s.mux.Handle("GET /api/dnssec", handlers.HandleDnsSec())
	s.mux.Handle("GET /api/firewall", handlers.HandleFirewall())
	s.mux.Handle("GET /api/get-ip", handlers.HandleGetIP(s.checks.IpAddress))
	s.mux.Handle("GET /api/headers", handlers.HandleGetHeaders(s.checks.Headers))
	s.mux.Handle("GET /api/hsts", handlers.HandleHsts())
	s.mux.Handle("GET /api/http-security", handlers.HandleHttpSecurity())
	s.mux.Handle("GET /api/legacy-rank", handlers.HandleLegacyRank(s.checks.LegacyRank))
	s.mux.Handle("GET /api/linked-pages", handlers.HandleGetLinks(s.checks.LinkedPages))
	s.mux.Handle("GET /api/ports", handlers.HandleGetPorts())
	s.mux.Handle("GET /api/quality", handlers.HandleGetQuality())
	s.mux.Handle("GET /api/rank", handlers.HandleGetRank(s.checks.Rank))
	s.mux.Handle("GET /api/redirects", handlers.HandleGetRedirects())
	s.mux.Handle("GET /api/social-tags", handlers.HandleGetSocialTags(s.checks.SocialTags))
	s.mux.Handle("GET /api/tls", handlers.HandleTLS(s.checks.Tls))
	s.mux.Handle("GET /api/trace-route", handlers.HandleTraceRoute())

	s.srv.Handler = s.CORS(s.mux)
}

func (s *Server) Run() error {
	s.routes()

	addr := fmt.Sprintf("%s:%s", s.conf.Host, s.conf.Port)
	log.Printf("Server started, listening on: %v\n", addr)
	s.srv.Addr = addr
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
