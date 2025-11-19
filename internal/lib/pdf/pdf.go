package pdf

import (
	"bytes"
	"web-server/internal/storage"

	"github.com/go-pdf/fpdf"
)

func GeneratePDF(links map[string]storage.LinkStatus) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	y := 10.0
	for link, status := range links {
		pdf.Text(10, y, link+": "+string(status))
		y += 6
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
