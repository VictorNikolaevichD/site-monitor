package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	addr   string
	log    *slog.Logger
	server *grpc.Server
}

func NewServer(addr string, log *slog.Logger, register func(*grpc.Server)) *Server {
	s := grpc.NewServer()
	register(s)

	return &Server{
		addr:   addr,
		log:    log,
		server: s,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("listen grpc %s: %w", s.addr, err)
	}

	s.log.Info("grpc server listening", "addr", s.addr)
	return s.server.Serve(lis)
}

func (s *Server) GracefulStop() {
	s.log.Info("grpc server shutting down")
	s.server.GracefulStop()
}
