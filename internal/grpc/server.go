package grpc

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/ViktorNikolaevichD/site-monitor/internal/grpc/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	addr   string
	logger *slog.Logger
	server *grpc.Server
}

func NewServer(addr string, logger *slog.Logger, register func(*grpc.Server)) *Server {
	s := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             30 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.ChainUnaryInterceptor(
			interceptor.UnaryRecover(logger),
			interceptor.UnaryLogging(logger),
		),
	)
	register(s)

	return &Server{
		addr:   addr,
		logger: logger,
		server: s,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("listen grpc %s: %w", s.addr, err)
	}

	s.logger.Info("grpc server listening", "addr", s.addr)
	return s.server.Serve(lis)
}

func (s *Server) GracefulStop() {
	s.logger.Info("grpc server shutting down")
	s.server.GracefulStop()
}
