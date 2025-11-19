package uploadpdf

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"web-server/internal/lib/pdf"
	"web-server/internal/storage"
)

type Request struct {
	LinksList []int64 `json:"links_list"`
}

type LinksLoader interface {
	LoadLinsksAndSatsuses(links_number int64) map[string]storage.LinkStatus
}

func New(log *slog.Logger, linksLoader LinksLoader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "links.uploadPDF.New"

		log := log.With("op", op)

		if r.Method != http.MethodPost {
			log.Error("invalid HTTP method", "method", r.Method)

			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		var req Request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request body", "error", err)

			w.WriteHeader(http.StatusBadRequest)

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		linksWithStat := make(map[string]storage.LinkStatus)

		for _, listNum := range req.LinksList {
			links := linksLoader.LoadLinsksAndSatsuses(listNum)
			for link := range links {
				linksWithStat[link] = links[link]
			}
		}

		pdfBytes, err := pdf.GeneratePDF(linksWithStat)
		if err != nil {
			log.Error("failed to generate pdf", "error", err)

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=\"report.pdf\"")

		if _, err := w.Write(pdfBytes); err != nil {
			log.Error("failed to write PDF", "error", err)

			w.WriteHeader(http.StatusInternalServerError)

			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
