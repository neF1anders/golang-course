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
	"repo-stat/processor/internal/adapter/collector"
	grpccontroller "repo-stat/processor/internal/controller/grpc"
	"repo-stat/processor/internal/usecase"
	pb "repo-stat/proto/processor"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "processor/config/config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting processor server...")
	log.Debug("debug messages are enabled")

	CollectorClient, err := collector.NewClient(cfg.Services.COLLECTOR, log)
	if err != nil {
		log.Info(fmt.Sprintf("unsuccessful connection to collector: %s", err))
	}

	pingUseCase := usecase.NewPing()
	fetchUseCase := usecase.NewFetch(CollectorClient) //??

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
