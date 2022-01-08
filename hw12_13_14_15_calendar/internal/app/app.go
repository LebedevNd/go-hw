package app

import (
	"context"

	"github.com/LebedevNd/go-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Storage interface {
	AddEvent(
		Title string,
		DateTimeStart string,
		DateTimeEnd string,
		Description string,
		UserID int,
		Delay int,
	) error
	UpdateEvent(
		ID int,
		Title string,
		DateTimeStart string,
		DateTimeEnd string,
		Description string,
		UserID int,
		Delay int,
	) error
	DeleteEvent(ID int) error
	GetEvents() ([]storage.Event, error)
	GetEventsByDay(day string) ([]storage.Event, error)
	GetEventsByWeek(day string) ([]storage.Event, error)
	GetEventsByMonth(day string) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
