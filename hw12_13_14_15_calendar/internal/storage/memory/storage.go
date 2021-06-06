package memorystorage

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) CreateEvent(e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := len(s.events)
	id++
	e.ID = strconv.Itoa(id)

	s.events[e.ID] = e
	return nil
}

func (s *Storage) UpdateEvent(e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[e.ID] = e
	return nil
}

func (s *Storage) DeleteEvent(e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.events, e.ID)
	return nil
}

func (s *Storage) GetEvents(ctx context.Context, startData time.Time, endData time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []storage.Event
	for _, event := range s.events {
		if event.StartData.Unix() >= startData.Unix() && event.EndData.Unix() <= endData.Unix() {
			events = append(events, event)
		}
	}
	return events, nil
}
