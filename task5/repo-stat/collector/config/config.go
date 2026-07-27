package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	"time"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-collector"`
}

type Services struct {
	Processor string `yaml:"processor" env:"PROCESSOR_ADDRESS" env-default:"processor:8083"`
	Subscribe string `yaml:"subscribe" env:"SUBSCRIBE_ADDRESS" env-default:"subscribe:8085"`
}
type Kafka struct {
	Brokers     string        `yaml:"brokers" env:"BROKERS" env-default:"kafka:9092"`
	Group       string        `yaml:"group" env:"GROUP" env-default:"collector-group"`
	InputTopic  string        `yaml:"intopic" env:"KAFKA_INPUT_TOPIC" env-default:"collect-commands"`
	OutputTopic string        `yaml:"outtopic" env:"KAFKA_OUTPUT_TOPIC" env-default:"collect-results"`
	Interval    time.Duration `yaml:"auto_update_interval" env:"AUTO_UPDATE_INTERVAL" env-default:"15s"`
}
type Config struct {
	App      App               `yaml:"app"`
	Services Services          `yaml:"services"`
	GRPC     grpcserver.Config `yaml:"grpc"`
	Kafka    Kafka             `yaml:"kafka"`
	Logger   logger.Config     `yaml:"logger"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
