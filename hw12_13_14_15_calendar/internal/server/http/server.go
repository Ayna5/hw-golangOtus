package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	Address string
	server  *http.Server
	logger  logger.Logger
}

type Application interface {
	// TODO
}

func NewServer(host, port string, log logger.Logger) *Server {
	return &Server{
		Address: net.JoinHostPort(host, port),
		logger:  log,
	}
}

func (s *Server) Start(ctx context.Context) error {
	router := mux.NewRouter()
	router.HandleFunc("/hello", s.helloWorld).Methods("GET")
	router.Use(s.loggingMiddleware)

	s.server = &http.Server{
		Addr:         s.Address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "start server error")
	}

	select {
	case <-ctx.Done():
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return errors.New("server is nil")
	}
	if err := s.server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "stop server error")
	}
	return nil
}

func (s *Server) helloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}
