package storage

type parsedData struct {
	NextId  int64       `json:"next_id"`
	Batches []LinkBatch `json:"batches"`
}

type LinkBatch struct {
	ID       int64
	Links    []string
	Statuses []string
}
