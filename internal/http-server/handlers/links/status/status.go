package status

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Request struct {
	Links []string `json:"links"`
}

type Response struct {
	Statuses []string `json:"statuses"`
	BachID   int64    `json:"batch_id"`
}

type LinksSaver interface {
	SaveLinksAndStatus(links []string, statuses []string) (int64, error)
}

type StatusGetter interface {
	GetStatuses(links []string) ([]string, error)
}

func New(log *slog.Logger, linksSaver LinksSaver, statusGetter StatusGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "status.links.Status.New"

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

		statuses, err := statusGetter.GetStatuses(req.Links)
		if err != nil {
			log.Error("failed to get statuses", "error", err)

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		batchID, err := linksSaver.SaveLinksAndStatus(req.Links, statuses)
		if err != nil {
			log.Error("failed to save links and statuses", "error", err)

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		log.Info("batch saved", slog.Int64("id", batchID))

		resp := Response{
			Statuses: statuses,
			BachID:   batchID,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("failed to encode response", "error", err)

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
