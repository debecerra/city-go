package http

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/debecerra/city-go/backend/internal/http/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router *chi.Mux
	http   *http.Server
	logger *slog.Logger
}

func NewServer(port string) *Server {
	s := &Server{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
	s.router = s.setupRouter()
	s.http = &http.Server{
		Addr:         ":" + port,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return s
}

func (s *Server) setupRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Routes
	r.Get("/health", handlers.Health)
	r.Route("/v1", func(r chi.Router) {
		r.Post("/recommend", handlers.Recommend)
	})

	return r
}

func (s *Server) Run() {
	go func() {
		s.logger.Info("server starting", "addr", s.http.Addr)
		if err := s.http.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.logger.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := s.http.Shutdown(ctx); err != nil {
		s.logger.Error("shutdown error", "err", err)
	}
	s.logger.Info("server stopped")
}
