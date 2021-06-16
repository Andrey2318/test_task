package server

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test_task/internal/config"
	"time"
)

type Server struct {
	config      *config.Config
	gRPCServer  *GRPCServer
	proxyServer *ProxyServer
}

func New(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (s *Server) Run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan error, 1)
	go func() {
		s.gRPCServer = &GRPCServer{
			server: s,
			Server: grpc.NewServer(
				grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
					grpc_recovery.StreamServerInterceptor(),
					grpc_ctxtags.StreamServerInterceptor(),
					grpc_opentracing.StreamServerInterceptor(),
					grpc_prometheus.StreamServerInterceptor,
				)),
				grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
					grpc_recovery.UnaryServerInterceptor(),
					grpc_ctxtags.UnaryServerInterceptor(),
					grpc_opentracing.UnaryServerInterceptor(),
					grpc_prometheus.UnaryServerInterceptor,
				))),
		}
		if err := s.gRPCServer.Run(); err != nil {
			ch <- err
		}
	}()

	go func() {
		s.proxyServer = &ProxyServer{
			server: s,
			Server: &http.Server{},
		}
		if err := s.proxyServer.Run(ctx); err != nil {
			ch <- err
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	select {
	case err := <-ch:
		return err
	case <-interrupt:
	}

	timeout, CancelFunc := context.WithTimeout(ctx, time.Second*10)
	defer CancelFunc()

	group := &errgroup.Group{}
	group.Go(func() error {
		if err := s.proxyServer.Shutdown(timeout); err != nil {
			return err
		}
		return nil
	})

	group.Go(func() error {
		s.gRPCServer.GracefulStop()
		return nil
	})

	return nil
}
