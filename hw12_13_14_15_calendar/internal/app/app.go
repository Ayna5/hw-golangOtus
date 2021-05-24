package app

import (
	"fmt"
	"time"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Storage interface {
	CreateEvent(e storage.Event) error
	UpdateEvent(e storage.Event) error
	DeleteEvent(e storage.Event) error
	GetEvents(startData, endData time.Time) ([]*storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(e storage.Event) error {
	if err := a.storage.CreateEvent(e); err != nil {
		a.logger.Error("cannot create event error")
		return fmt.Errorf("cannot create event: %w", err)
	}
	return nil
}

func (a *App) UpdateEvent(e storage.Event) error {
	if err := a.storage.UpdateEvent(e); err != nil {
		a.logger.Error("cannot update event error")
		return fmt.Errorf("cannot update event: %w", err)
	}
	return nil
}

func (a *App) DeleteEvent(e storage.Event) error {
	if err := a.storage.DeleteEvent(e); err != nil {
		a.logger.Error("cannot delete event error")
		return fmt.Errorf("cannot delete event: %w", err)
	}
	return nil
}

func (a *App) GetEvents(startData, endData time.Time) ([]*storage.Event, error) {
	events, err := a.storage.GetEvents(startData, endData)
	if err != nil {
		a.logger.Error("cannot get events error")
		return nil, fmt.Errorf("cannot get events: %w", err)
	}
	return events, nil
}
