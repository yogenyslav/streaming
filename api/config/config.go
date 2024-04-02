package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/yogenyslav/storage/postgres"
)

type Config struct {
	Server       ServerConfig       `yaml:"server"`
	Postgres     postgres.Config    `yaml:"postgres"`
	FrameService FrameServiceConfig `yaml:"frameService"`
}

type ServerConfig struct {
	Port        int    `yaml:"port"`
	LogLevel    string `yaml:"logLevel"`
	CorsOrigins string `yaml:"corsOrigins"`
}

type FrameServiceConfig struct {
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
