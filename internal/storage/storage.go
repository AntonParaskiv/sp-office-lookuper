package storage

import (
	"sync"

	"sp-office-lookuper/internal/app"
)

type Storage struct {
	offices map[int64]int64
	mu      *sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		offices: map[int64]int64{},
		mu:      &sync.Mutex{},
	}
}

func (s *Storage) SetOfficeSortPoint(officeID, sortPointID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.offices[officeID] = sortPointID
}

func (s *Storage) GetSortPoint(officeID int64) (int64, error) {
	s.mu.Lock()
	sortPointID, ok := s.offices[officeID]
	s.mu.Unlock()

	if !ok {
		return 0, app.ErrOfficeNotFound
	}
	return sortPointID, nil
}
