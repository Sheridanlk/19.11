package storage

import (
	"encoding/json"
	"os"
	"sync"
)

type Storage struct {
	path string

	mu      sync.RWMutex
	NextID  int64
	batches map[int64]*LinkBatch
}

func New(path string) *Storage {
	return &Storage{
		path:    path,
		batches: make(map[int64]*LinkBatch),
		NextID:  1,
	}
}

func (s *Storage) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var parsedData parsedData
	if err := json.Unmarshal(data, &parsedData); err != nil {
		return err
	}
	s.NextID = parsedData.NextId
	s.batches = make(map[int64]*LinkBatch, len(parsedData.Batches))
	for _, batch := range parsedData.Batches {
		s.batches[batch.ID] = &LinkBatch{
			ID:    batch.ID,
			Links: batch.Links,
		}
	}
	return nil
}

func (s *Storage) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	parsedData := parsedData{
		NextId:  s.NextID,
		Batches: make([]LinkBatch, 0, len(s.batches)),
	}
	for _, batch := range s.batches {
		parsedData.Batches = append(parsedData.Batches, *batch)
	}

	data, err := json.MarshalIndent(parsedData, "", "  ")
	if err != nil {
		return err
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func (s *Storage) SaveLinksAndStatuses(linksAndStatus map[string]LinkStatus) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.NextID
	s.NextID++

	batch := &LinkBatch{
		ID:    id,
		Links: linksAndStatus,
	}

	s.batches[id] = batch

	return id, nil
}

func (s *Storage) LoadLinsksAndSatsuses(link_num int64) map[string]LinkStatus {
	s.mu.Lock()
	defer s.mu.Unlock()

	batch := s.batches[link_num]

	return batch.Links
}
