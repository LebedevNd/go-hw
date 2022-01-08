package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/LebedevNd/go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/golang-module/carbon"

	// comment for lint.
	_ "github.com/jackc/pgx/stdlib"
)

type Storage struct {
	username string
	password string
	host     string
	port     int
	dbname   string
	db       *sql.DB
}

func (s *Storage) CheckDateInterval(id int, dateFrom string, dateTo string) error {
	whereQuery := fmt.Sprintf("where date_time_start between "+
		"%s and %s or date_time_end between %s and %s",
		dateFrom, dateTo, dateFrom, dateTo)

	if id > 0 {
		whereQuery += fmt.Sprintf(" and id != %d", id)
	}

	events, err := s.RequestEvents(whereQuery)
	if err != nil {
		return err
	}

	if len(events) > 0 {
		return storage.BusyDateError{
			Title:    events[0].Title,
			DateFrom: events[0].DateTimeStart,
			DateTo:   events[0].DateTimeEnd,
		}
	}

	return nil
}

func New(username string, password string, host string, port int, dbname string) *Storage {
	return &Storage{
		username,
		password,
		host,
		port,
		dbname,
		nil,
	}
}

func (s *Storage) Connect(
	ctx context.Context,
) error {
	dsn := "user=" + s.username +
		" dbname=" + s.dbname +
		" sslmode=disable password=" + s.password +
		" host=" + s.host +
		" port=" + strconv.Itoa(s.port)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Storage) Close() error {
	if s.db != nil {
		err := s.db.Close()
		return err
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
	err := s.CheckDateInterval(0, dateTimeStart, dateTimeEnd)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = s.Connect(ctx)
	if err != nil {
		return err
	}
	defer func(s *Storage) {
		err := s.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(s)

	query := "insert into " +
		"events(title, date_time_start, date_time_end, description, user_id, delay)" +
		" values($1, $2, $3, $4, $5, $6)"
	_, err = s.db.Exec(query, title, dateTimeStart, dateTimeEnd, description, userID, delay)

	return err
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
	err := s.CheckDateInterval(id, dateTimeStart, dateTimeEnd)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = s.Connect(ctx)
	if err != nil {
		return err
	}
	defer func(s *Storage) {
		err := s.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(s)

	query := "update events set " +
		"title=$1, date_time_start=$2, date_time_end=$3, description=$4, user_id=$5, delay=$6 where id=$7"
	_, err = s.db.Exec(query, title, dateTimeStart, dateTimeEnd, description, userID, delay, id)
	return err
}

func (s *Storage) DeleteEvent(id int) error {
	ctx := context.Background()
	err := s.Connect(ctx)
	if err != nil {
		return err
	}
	defer func(s *Storage) {
		err := s.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(s)

	query := "delete from events where id=$1"
	_, err = s.db.Exec(query, id)
	return err
}

func (s *Storage) GetEvents() ([]storage.Event, error) {
	events, err := s.RequestEvents("")
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) RequestEvents(whereQuery string) ([]storage.Event, error) {
	ctx := context.Background()
	err := s.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer func(s *Storage) {
		err := s.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(s)

	query := "select * from events " + whereQuery
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var id, userID, delay int
		var title, dateTimeStart, dateTimeEnd, description string
		if err := rows.Scan(&id, &title, &dateTimeStart, &dateTimeEnd, &description, &userID, &delay); err != nil {
			fmt.Println(err)
		}
		events = append(events, storage.Event{
			ID:            id,
			Description:   description,
			Delay:         delay,
			DateTimeStart: dateTimeStart,
			DateTimeEnd:   dateTimeEnd,
			UserID:        userID,
			Title:         title,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
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

	whereQuery := fmt.Sprintf(" and where date_time_start >= %s and date_time_start < %s", dateStart, dateEnd)
	return s.RequestEvents(whereQuery)
}
