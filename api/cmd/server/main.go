package main

import (
	"streaming/api/config"
	"streaming/api/internal/server"

	"github.com/rs/zerolog"
	"github.com/yogenyslav/pkg/loctime"
)

func main() {
	cfg := config.MustNew("./config/config.yaml")
	level, err := zerolog.ParseLevel(cfg.Server.LogLevel)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(level)
	loctime.SetLocation(loctime.MoscowLocation)
	srv := server.New(cfg)
	srv.Run()
}
