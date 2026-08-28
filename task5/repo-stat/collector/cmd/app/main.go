package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/collector/config"
	"repo-stat/platform/logger"

	"repo-stat/collector/internal/adapter/broker"
	"repo-stat/collector/internal/adapter/github"
	"repo-stat/collector/internal/adapter/grpc"
	"repo-stat/collector/internal/controller/consumer"
	"repo-stat/collector/internal/controller/producer"

	"repo-stat/collector/internal/usecase"
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

	githubClient := github.NewClient(log)
	subscribeClient, err := grpc.NewClient(cfg.Services.Subscribe, log)
	if err != nil {
		log.Info(fmt.Sprintf("unsuccessful connection to subscribe: %s", err))
	}

	kafkaProducer := broker.NewProducer([]string{cfg.Kafka.Brokers}, cfg.Kafka.OutputTopic, log)
	defer func() {
		err := kafkaProducer.Close()
		if err != nil {
			log.Error("failed to close producer", "error", err)
		}
	}()
	log.Debug("adapters are loaded successfully")

	kafkaPublisher := producer.NewResultPublisher(kafkaProducer, cfg.Kafka.OutputTopic)
	getSubInfoUseCase := usecase.NewGetSubInfoUseCase(githubClient, subscribeClient)
	getAndPublishUseCase := usecase.NewGetAndPublishUseCase(getSubInfoUseCase, kafkaPublisher)
	log.Debug("usecases are loaded successfully")

	updateScheduler := consumer.NewScheduler(log, getAndPublishUseCase, cfg.Kafka.Interval)
	messageHandler := consumer.NewOrderMessageHandler(getAndPublishUseCase)
	log.Debug("controllers are loaded successfully")

	kafkaConsumer, err := broker.NewConsumer([]string{cfg.Kafka.Brokers}, cfg.Kafka.Group, cfg.Kafka.InputTopic, messageHandler, log)
	if err != nil {
		log.Error("failed to initialize consumer", "error", err)
		return err
	}
	defer func() {
		err := kafkaConsumer.Stop()
		if err != nil {
			log.Error("failed to close consumer", "error", err)
		}
	}()
	if err := kafkaConsumer.Start(ctx); err != nil {
		log.Error("failed to start kafka consumer", "error", err)
		return err
	}
	updateScheduler.Start(ctx)
	log.Info("collector is up successfully")
	<-ctx.Done()
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
