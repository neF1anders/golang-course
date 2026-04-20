package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	"repo-stat/subscribe/config"
	"repo-stat/subscribe/internal/usecase"

	subscribepb "repo-stat/proto/subscribe"
	adapter "repo-stat/subscribe/internal/adapter/github"
	"repo-stat/subscribe/internal/adapter/postgres"
	grpccontroller "repo-stat/subscribe/internal/controller/grpc"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "subscribe/config/config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting subscribe server...")
	log.Debug("debug messages are enabled")

	githubClient := adapter.NewClient(log)
	pingUseCase := usecase.NewPingerUseCase(githubClient)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host,
		cfg.Database.Port, cfg.Database.DBname, cfg.Database.SSLmode)
	m, err := migrate.New("file://db/migrations", dsn)
	if err != nil {
		return fmt.Errorf("create migrate.New: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("up migrate: %w", err)
	}
	pgclient, err := postgres.NewPostgresClient(ctx, dsn)
	if err != nil {
		return fmt.Errorf("create PostgreSQL client: %w", err)
	}
	db := usecase.NewDBUseCase(pgclient)

	server := grpccontroller.NewSubscribeServer(log, db, pingUseCase)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	subscribepb.RegisterSubscribeServer(srv.GRPC(), server)
	log.Info(fmt.Sprintf("Subscribe gRPC server listening on :%v", cfg.GRPC.Address))
	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	if err := run(ctx); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Printf("launching server error: %s\n", err)
		}
		cancel()
		os.Exit(1)
	}
}
