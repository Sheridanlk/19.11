package main

import (
	"web-server/internal/config"
	"web-server/internal/logger"
	"web-server/internal/service"
)

func main() {
	cfg := config.Load()
	log := logger.SetupLogger(cfg.Env)

	service := service.New(log, cfg.Data_path, cfg.HTTPServer.Address, cfg.HTTPServer.Timeout, cfg.HTTPServer.Timeout, cfg.HTTPServer.IdleTimeout)
	go func() {
		service.Start()
	}()
}
