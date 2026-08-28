package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-processor"`
}

type Services struct {
	API       string `yaml:"api" env:"API_ADDRESS" env-default:"localhost:8080"`
	COLLECTOR string `yaml:"collector" env:"COLLECTOR_ADDRESS" env-default:"localhost:8083"`
}

type Kafka struct {
	Brokers     string `yaml:"brokers" env:"BROKERS" env-default:"kafka:9092"`
	Group       string `yaml:"group" env:"GROUP" env-default:"collector-group"`
	InputTopic  string `yaml:"intopic" env:"KAFKA_INPUT_TOPIC" env-default:"collect-resullts"`
	OutputTopic string `yaml:"outtopic" env:"KAFKA_OUTPUT_TOPIC" env-default:"collect-commands"`
}

type DB struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBname   string `yaml:"dbname"`
	SSLmode  string `yaml:"sslmode"`
}

type Config struct {
	App      App               `yaml:"app"`
	Services Services          `yaml:"services"`
	GRPC     grpcserver.Config `yaml:"grpc"`
	Kafka    Kafka             `yaml:"kafka"`
	Database DB                `yaml:"database"`
	Logger   logger.Config     `yaml:"logger"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
