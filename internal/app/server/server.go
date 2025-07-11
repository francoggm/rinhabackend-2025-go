package server

import (
	"fmt"
	"francoggm/rinhabackend-2025-go/internal/app/server/handlers"
	"francoggm/rinhabackend-2025-go/internal/config"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Server struct {
	cfg      *config.Config
	router   *chi.Mux
	handlers *handlers.Handlers
}

func NewServer(cfg *config.Config, db *pgxpool.Pool, paymentEventsCh chan any) *Server {
	srv := &Server{
		cfg:      cfg,
		router:   chi.NewRouter(),
		handlers: handlers.NewHandlers(cfg, db, paymentEventsCh),
	}

	srv.registerRoutes()
	return srv
}

func (s *Server) registerRoutes() {
	s.router.Post("/payments", s.handlers.ProcessPayment)
	s.router.Get("/payments-summary", s.handlers.GetPaymentsSummary)
}

func (s *Server) Run() error {
	zap.L().Info("Starting server", zap.String("port", s.cfg.Server.Port))
	return http.ListenAndServe(fmt.Sprintf(":%s", s.cfg.Server.Port), s.router)
}
