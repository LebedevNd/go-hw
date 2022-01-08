package storage

import "fmt"

type BusyDateError struct {
	Title    string
	DateFrom string
	DateTo   string
}

func (e BusyDateError) Error() string {
	return fmt.Sprintf("DateTime is busy with %s event (from %s to %s)", e.Title, e.DateFrom, e.DateTo)
}

type Event struct {
	ID            int
	Title         string
	DateTimeStart string
	DateTimeEnd   string
	Description   string
	UserID        int
	Delay         int
}
