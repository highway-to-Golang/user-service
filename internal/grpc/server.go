package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/highway-to-Golang/user-service/api/proto/gen/go/user"
	"github.com/highway-to-Golang/user-service/config"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	addr       string
}

func NewServer(cfg config.Config, userHandler *UserHandler) *Server {
	grpcServer := grpc.NewServer()

	user.RegisterUserServiceServer(grpcServer, userHandler)

	addr := fmt.Sprintf("%s:%s", cfg.GRPC.Host, cfg.GRPC.Port)

	return &Server{
		grpcServer: grpcServer,
		addr:       addr,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.addr, err)
	}

	slog.Info("starting gRPC server", "address", s.addr)

	if err := s.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("shutting down gRPC server")

	done := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		s.grpcServer.Stop()
		return ctx.Err()
	case <-done:
		slog.Info("gRPC server stopped")
		return nil
	}
}
