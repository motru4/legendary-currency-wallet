package grpc

import (
	"log/slog"
	"net"

	"github.com/motru4/legendary-currency-wallet/proto/exchange"
	"github.com/motru4/legendary-currency-wallet/services/exchanger/internal/repositories"

	grpcLib "google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpcLib.Server
	log        *slog.Logger
}

func NewServer(repo repositories.ExchangeRatesRepository, log *slog.Logger) *Server {
	grpcSrv := grpcLib.NewServer()

	exchangeService := NewExchangeServiceServer(repo, log)
	exchange.RegisterExchangeServiceServer(grpcSrv, exchangeService)

	return &Server{
		grpcServer: grpcSrv,
		log:        log.With("component", "grpc_server"),
	}
}

func (s *Server) Serve(listener net.Listener) error {
	s.log.Info("starting gRPC server", "addr", listener.Addr().String())
	return s.grpcServer.Serve(listener)
}

func (s *Server) GracefulStop() {
	s.log.Info("stopping gRPC server gracefully")
	s.grpcServer.GracefulStop()
}
