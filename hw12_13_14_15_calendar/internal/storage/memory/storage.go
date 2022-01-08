package memorystorage

import (
	"errors"
	"sync"

	"github.com/LebedevNd/go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/golang-module/carbon"
)

type Storage struct {
	events []storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	var events []storage.Event
	return &Storage{
		events,
		sync.RWMutex{},
	}
}

func (s *Storage) CheckDateInterval(id int, dateFrom string, dateTo string) error {
	for _, event := range s.events {
		if (dateFrom >= event.DateTimeStart && dateFrom <= event.DateTimeEnd ||
			dateTo >= event.DateTimeStart && dateTo <= event.DateTimeEnd) && !(id > 0 && event.ID == id) {
			return storage.BusyDateError{
				Title:    event.Title,
				DateFrom: event.DateTimeStart,
				DateTo:   event.DateTimeEnd,
			}
		}
	}
	return nil
}

func (s *Storage) AddEvent(
	title string,
	dateTimeStart string,
	dateTimeEnd string,
	description string,
	userID int,
	delay int,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.CheckDateInterval(0, dateTimeStart, dateTimeEnd)
	if err != nil {
		return err
	}

	event := storage.Event{
		ID:            s.GenerateID(),
		Title:         title,
		DateTimeStart: dateTimeStart,
		DateTimeEnd:   dateTimeEnd,
		Description:   description,
		UserID:        userID,
		Delay:         delay,
	}
	s.events = append(s.events, event)
	// TODO errors/validation
	return nil
}

func (s *Storage) GenerateID() int {
	if len(s.events) == 0 {
		return 1
	}

	lastEvent := s.events[len(s.events)-1]
	return lastEvent.ID + 1
}

func (s *Storage) DeleteEvent(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, event := range s.events {
		if event.ID == id {
			s.events = append(s.events[:i], s.events[i+1:]...)
			return nil
		}
	}
	return errors.New("event not found ")
}

func (s *Storage) UpdateEvent(
	id int,
	title string,
	dateTimeStart string,
	dateTimeEnd string,
	description string,
	userID int,
	delay int,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.CheckDateInterval(id, dateTimeStart, dateTimeEnd)
	if err != nil {
		return err
	}

	updatedEvent := storage.Event{
		ID:            s.GenerateID(),
		Title:         title,
		DateTimeStart: dateTimeStart,
		DateTimeEnd:   dateTimeEnd,
		Description:   description,
		UserID:        userID,
		Delay:         delay,
	}

	for i, event := range s.events {
		if event.ID == id {
			s.events[i] = updatedEvent
			return nil
		}
	}
	return errors.New("event not found ")
}

func (s *Storage) GetEvents() ([]storage.Event, error) {
	return s.events, nil
}

func (s *Storage) GetEventsByDay(day string) ([]storage.Event, error) {
	return s.GetEventsBy(day, "day")
}

func (s *Storage) GetEventsByWeek(day string) ([]storage.Event, error) {
	return s.GetEventsBy(day, "week")
}

func (s *Storage) GetEventsByMonth(day string) ([]storage.Event, error) {
	return s.GetEventsBy(day, "month")
}

func (s *Storage) GetEventsBy(day string, interval string) ([]storage.Event, error) {
	dayParsed := carbon.Parse(day)
	dateStart := dayParsed.Format("Y-m-d")

	var dateEnd string
	switch interval {
	case "day":
		dateEnd = dayParsed.AddDay().Format("Y-m-d")
	case "week":
		dateEnd = dayParsed.AddWeek().Format("Y-m-d")
	case "month":
		dateEnd = dayParsed.AddMonth().Format("Y-m-d")
	default:
		return nil, errors.New("wrong interval name ")
	}

	var events []storage.Event
	for _, event := range s.events {
		if event.DateTimeStart >= dateStart && event.DateTimeStart < dateEnd {
			events = append(events, event)
		}
	}

	return events, nil
}
