package config

type ServerConfig struct {
	Port     int    `yaml:"port"`
	LogLevel string `yaml:"log_level"`
}
