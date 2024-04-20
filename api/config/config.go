package config

import (
	"streaming/api/pkg/storage/minios3"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/yogenyslav/storage/postgres"
)

type Config struct {
	Server       ServerConfig    `yaml:"server"`
	Postgres     postgres.Config `yaml:"postgres"`
	Orchestrator ServiceConfig   `yaml:"orchestrator"`
	S3Config     minios3.Config  `yaml:"minio"`
}

type ServerConfig struct {
	Port        int    `yaml:"port"`
	LogLevel    string `yaml:"logLevel"`
	CorsOrigins string `yaml:"corsOrigins"`
}

type ServiceConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func MustNew(path string) *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic(err)
	}
	return cfg
}
