package grpc

import (
	"fmt"
	"net"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/app"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/logger"
	calendar_pb "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/pkg/calendar"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Server struct {
	calendar_pb.UnimplementedCalendarServer
	lis    net.Listener
	l      *logger.Logger
	server *grpc.Server
	api    *app.App
}

func NewServer(l *logger.Logger, api *app.App, address string) (*Server, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("start listen error: %w", err)
	}

	server := grpc.NewServer()
	srv := &Server{
		lis:    lis,
		l:      l,
		server: server,
		api:    api,
	}
	calendar_pb.RegisterCalendarServer(server, srv)

	return srv, nil
}

func (s *Server) Start() error {
	if err := s.server.Serve(s.lis); err != nil {
		return fmt.Errorf("start server error: %w", err)
	}
	return nil
}

func (s *Server) Stop() error {
	if s.server == nil {
		return errors.New("grpc server is nil") //nolint:errcheck,govet
	}

	s.server.GracefulStop()
	return nil
}
