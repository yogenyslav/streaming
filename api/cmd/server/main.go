package main

import (
	"streaming/api/config"
	"streaming/api/internal/server"

	"github.com/yogenyslav/logger"
)

func main() {
	cfg := config.MustNew("./config/config.yaml")
	logger.SetLevel(logger.ParseLevel(cfg.Server.LogLevel))

	app := server.New(cfg)
	app.Run()
}
