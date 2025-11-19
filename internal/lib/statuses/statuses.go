package statuses

import (
	"net/http"
	"strings"
	"time"
	"web-server/internal/storage"
)

var client = http.Client{
	Timeout: 5 * time.Second,
}

func GetStatus(link string) storage.LinkStatus {
	url := normalizeLink(link)

	resp, err := client.Get(url)
	if err != nil {
		return storage.StatusNotAvailable
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return storage.StatusAvailable
	} else {
		return storage.StatusNotAvailable
	}
}

func normalizeLink(link string) string {
	if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
		return link
	}
	// по умолчанию считаем https
	return "https://" + link
}
