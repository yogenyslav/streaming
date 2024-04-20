package main

import (
	"streaming/orchestrator/config"
	"streaming/orchestrator/internal/server"
)

func main() {
	cfg := config.MustNew("./config/config.yaml")

	app := server.New(cfg)
	app.Run()
}
