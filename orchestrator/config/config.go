package config

import (
	"streaming/orchestrator/pkg/infrastructure/kafka"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server ServerConfig  `yaml:"server"`
	Kafka  *kafka.Config `yaml:"kafka"`
}

type ServerConfig struct {
	Port     int    `yaml:"port"`
	LogLevel string `yaml:"logLevel"`
}

func MustNew(path string) *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic(err)
	}
	return cfg
}
