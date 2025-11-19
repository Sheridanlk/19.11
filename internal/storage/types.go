package storage

type LinkStatus string

const (
	StatusAvailable    LinkStatus = "available"
	StatusNotAvailable LinkStatus = "not available"
)

type parsedData struct {
	NextId  int64       `json:"next_id"`
	Batches []LinkBatch `json:"batches"`
}

type LinkBatch struct {
	ID    int64                 `json:"id"`
	Links map[string]LinkStatus `json:"links"`
}
