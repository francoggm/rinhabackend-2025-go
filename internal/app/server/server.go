package server

import (
	"fmt"
	"francoggm/rinhabackend-2025-go/internal/app/server/handlers"
	"francoggm/rinhabackend-2025-go/internal/config"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	cfg      *config.Config
	router   *chi.Mux
	handlers *handlers.Handlers
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:      cfg,
		router:   chi.NewRouter(),
		handlers: handlers.NewHandlers(),
	}
}

func (s *Server) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%s", s.cfg.Port), s.router)
}
