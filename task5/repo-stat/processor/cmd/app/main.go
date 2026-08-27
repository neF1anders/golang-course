package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	"repo-stat/processor/config"
	"repo-stat/processor/internal/adapter/broker"
	"repo-stat/processor/internal/adapter/postgres"
	grpccontroller "repo-stat/processor/internal/controller/grpc"
	"repo-stat/processor/internal/usecase"
	pb "repo-stat/proto/processor"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "processor/config/config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting processor server...")
	log.Debug("debug messages are enabled")

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
	pgclient, err := postgres.NewPostgresClient(ctx, dsn) //infra
	if err != nil {
		return fmt.Errorf("create PostgreSQL client: %w", err)
	}

	kafkaProducer := broker.NewProducer([]string{cfg.Kafka.Brokers}, cfg.Kafka.OutputTopic, log)
	defer kafkaProducer.Close()

	dbUseCase := usecase.NewDBUseCase(pgclient)
	pingUseCase := usecase.NewPing()
	fetchUseCase := usecase.NewFetchUC(dbUseCase, kafkaProducer)

	kafkaConsumer, err := broker.NewConsumer(
		[]string{cfg.Kafka.Brokers},
		cfg.Kafka.Group,
		cfg.Kafka.InputTopic,
		dbUseCase,
		log,
	)
	if err != nil {
		return fmt.Errorf("kafka consumer: %w", err)
	}
	defer kafkaConsumer.Stop()
	if err := kafkaConsumer.Start(ctx); err != nil {
		return err
	}

	processorHandler := grpccontroller.NewProcessorServer(log, fetchUseCase, pingUseCase)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("init grpc server: %w", err)
	}

	pb.RegisterProcessorServer(srv.GRPC(), processorHandler)
	log.Info(fmt.Sprintf("Processor gRPC server listening on :%v", cfg.GRPC.Address))
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
