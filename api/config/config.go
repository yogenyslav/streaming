package config

import (
	srvconf "streaming/api/internal/server/config"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/minios3"
	"github.com/yogenyslav/pkg/storage/postgres"
)

type Config struct {
	Server       *srvconf.ServerConfig `yaml:"server"`
	Postgres     *postgres.Config      `yaml:"postgres"`
	S3           *minios3.Config       `yaml:"s3"`
	Tracing      *tracing.Config       `yaml:"tracing"`
	Prometheus   *ServiceConfig        `yaml:"prometheus"`
	Orchestrator *ServiceConfig        `yaml:"orchestrator"`
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
