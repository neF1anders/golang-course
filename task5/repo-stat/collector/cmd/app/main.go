package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/collector/config"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"

	"repo-stat/collector/internal/adapter/github"
	"repo-stat/collector/internal/adapter/grpc"

	deliverygrpc "repo-stat/collector/internal/controller/grpc"
	"repo-stat/collector/internal/usecase"

	pb "repo-stat/proto/collector"
)

func run(ctx context.Context) error {
	// config
	var configPath string
	flag.StringVar(&configPath, "config", "collector/config/config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting collector server...")
	log.Debug("debug messages are enabled")

	githubClient := github.NewClient()
	subscribeClient, err := grpc.NewClient(cfg.Services.Subscribe, log)
	if err != nil {
		log.Info(fmt.Sprintf("unsuccessful connection to subscribe: %s", err))
	}

	getRepoInfoUseCase := usecase.NewGetRepoInfoUseCase(githubClient)
	getSubInfoUseCase := usecase.NewGetSubInfoUseCase(githubClient, subscribeClient)

	collectorServer := deliverygrpc.NewCollectorServer(getRepoInfoUseCase, getSubInfoUseCase)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("init grpc server: %w", err)
	}
	pb.RegisterCollectorServer(srv.GRPC(), collectorServer)

	log.Info(fmt.Sprintf("Collector gRPC server listening on :%v", cfg.GRPC.Address))
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
