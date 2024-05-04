package config

type ServerConfig struct {
	Port        int    `yaml:"port"`
	LogLevel    string `yaml:"log_level"`
	BodyLimit   int    `yaml:"body_limit"`
	CorsOrigins string `yaml:"cors_origins"`
}
