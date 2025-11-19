package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"web-server/internal/config"
	"web-server/internal/logger"
	"web-server/internal/service"
)

func main() {
	cfg := config.Load()
	log := logger.SetupLogger(cfg.Env)

	log.Info("Starting service", slog.String("env", cfg.Env))

	service := service.New(log, cfg.Data_path, cfg.HTTPServer.Address, cfg.HTTPServer.Timeout, cfg.HTTPServer.Timeout, cfg.HTTPServer.IdleTimeout)

	go func() {
		service.Start()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	service.Stop()
	log.Info("Service stopped")
}
