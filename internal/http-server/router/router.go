package router

import (
	"log/slog"
	"net/http"
	"web-server/internal/storage"
)

func Setup(log *slog.Logger, storage *storage.Storage) *http.ServeMux {
	rout := http.NewServeMux()
	// TODO: add routes

	return rout
}
