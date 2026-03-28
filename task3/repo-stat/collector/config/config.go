package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/httpserver"
	"repo-stat/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-collector"`
}

type Services struct {
	Processor string `yaml:"processor" env:"PROCESSOR_ADDRESS" env-default:"localhost:1234"`
}

type Config struct {
	App      App               `yaml:"app"`
	Services Services          `yaml:"services"`
	HTTP     httpserver.Config `yaml:"http"`
	Logger   logger.Config     `yaml:"logger"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
