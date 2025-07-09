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
	srv := &Server{
		cfg:      cfg,
		router:   chi.NewRouter(),
		handlers: handlers.NewHandlers(cfg),
	}

	srv.registerRoutes()
	return srv
}

func (s *Server) registerRoutes() {
	s.router.Post("/payments", s.handlers.ProcessPayment)
	s.router.Get("/payments-sumarry", s.handlers.GetPaymentsSummary)
}

func (s *Server) Run() error {
	fmt.Println("Starting server on port:", s.cfg.Port)
	return http.ListenAndServe(fmt.Sprintf(":%s", s.cfg.Port), s.router)
}
