package router

import (
	"log/slog"
	"net/http"
	"web-server/internal/http-server/handlers/links/status"
	"web-server/internal/http-server/handlers/links/uploadpdf"
	"web-server/internal/storage"
)

func Setup(log *slog.Logger, storage *storage.Storage) *http.ServeMux {
	rout := http.NewServeMux()
	rout.HandleFunc("/links", status.New(log, storage))
	rout.HandleFunc("/links/pdf", uploadpdf.New(log, storage))
	// TODO: add routes

	return rout
}
