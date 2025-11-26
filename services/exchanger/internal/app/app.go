package app

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/jmoiron/sqlx"
	commonpg "github.com/motru4/legendary-currency-wallet/pkg/db/postgres"
	commonlogger "github.com/motru4/legendary-currency-wallet/pkg/logger"
	"github.com/motru4/legendary-currency-wallet/services/exchanger/internal/config"
	repoiface "github.com/motru4/legendary-currency-wallet/services/exchanger/internal/repositories"
	repopg "github.com/motru4/legendary-currency-wallet/services/exchanger/internal/repositories/postgres"
	grpcTransport "github.com/motru4/legendary-currency-wallet/services/exchanger/internal/transport/grpc"
)

type App struct {
	cfg *config.Config

	log *slog.Logger
	db  *sqlx.DB

	grpcServer *grpcTransport.Server
	listener   net.Listener
}

func New(cfg *config.Config) (*App, error) {
	// 1. Logger
	log := commonlogger.New(cfg.Log.Level).With("service", "exchanger")

	log.Info("initializing exchanger service")

	// 2. Connecting to PostgreSQL
	db, err := commonpg.NewPostgres(commonpg.PostgresConfig{
		URL:             cfg.DB.URL(),
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		log.Error("failed to connect to postgres", "err", err)
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	log.Info("connected to postgres")

	// 3. Exchange rate repository
	var repo repoiface.ExchangeRatesRepository = repopg.NewExchangeRatesPostgresRepository(db, log)

	// 4. gRPC server
	grpcSrv := grpcTransport.NewServer(repo, log)

	// 5. Listener
	addr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("failed to listen", "addr", addr, "err", err)
		_ = db.Close()
		return nil, fmt.Errorf("listen %s: %w", addr, err)
	}

	log.Info("gRPC listener created", "addr", addr)

	return &App{
		cfg:        cfg,
		log:        log,
		db:         db,
		grpcServer: grpcSrv,
		listener:   lis,
	}, nil
}

func (a *App) Run() error {
	a.log.Info("starting gRPC server",
		"host", a.cfg.GRPC.Host,
		"port", a.cfg.GRPC.Port,
	)

	if err := a.grpcServer.Serve(a.listener); err != nil {
		a.log.Error("gRPC server stopped with error", "err", err)
		return err
	}

	a.log.Info("gRPC server stopped")
	return nil
}

// graceful shutdown
func (a *App) Stop() error {
	a.log.Info("stopping application")

	// 1. Stop the gRPC server
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}

	// 2. Close the database
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			a.log.Error("failed to close db", "err", err)
			return err
		}
	}

	a.log.Info("application stopped gracefully")
	return nil
}
