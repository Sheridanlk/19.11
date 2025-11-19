package service

import (
	"context"
	"log/slog"
	"net/http"
	"time"
	"web-server/internal/http-server/router"
	"web-server/internal/storage"
)

type Service struct {
	Server  *http.Server
	Storage *storage.Storage
}

func New(log *slog.Logger, path string, addres string, readTimeout time.Duration, writeTimeout time.Duration, idleTimeout time.Duration) *Service {
	storage := storage.New(path)

	router := router.Setup(log, storage)

	server := &http.Server{
		Addr:         addres,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	return &Service{
		Server:  server,
		Storage: storage,
	}
}

func (s *Service) Start() error {
	const op = "service.Start"

	s.Storage.Load()
	log := slog.With("op", op)
	log.Info("Starting server", "address", s.Server.Addr)

	return s.Server.ListenAndServe()
}

func (s *Service) Stop() {
	const op = "service.Stop"

	log := slog.With("op", op)
	log.Info("Stopping server and saving data")
	s.Server.Shutdown(context.Background())
	s.Storage.Save()
}
