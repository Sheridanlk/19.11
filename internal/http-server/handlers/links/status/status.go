package status

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"web-server/internal/lib/statuses"
	"web-server/internal/storage"
)

type Request struct {
	Links []string `json:"links"`
}

type Response struct {
	Links    map[string]storage.LinkStatus `json:"links"`
	LinksNum int64                         `json:"links_num"`
}

type LinksSaver interface {
	SaveLinksAndStatuses(map[string]storage.LinkStatus) (int64, error)
}

func New(log *slog.Logger, linksSaver LinksSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "links.Status.New"

		log := log.With("op", op)

		if r.Method != http.MethodPost {
			log.Error("invalid HTTP method", "method", r.Method)

			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		var req Request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("failed to decode request body", "error", err)

			w.WriteHeader(http.StatusBadRequest)

			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		linksWithStat := make(map[string]storage.LinkStatus)
		for _, link := range req.Links {
			linksWithStat[link] = statuses.GetStatus(link)
		}

		batchID, err := linksSaver.SaveLinksAndStatuses(linksWithStat)
		if err != nil {
			log.Error("failed to save links and statuses", "error", err)

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		log.Info("batch saved", slog.Int64("id", batchID))

		resp := Response{
			Links:    linksWithStat,
			LinksNum: batchID,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("failed to encode response", "error", err)

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
