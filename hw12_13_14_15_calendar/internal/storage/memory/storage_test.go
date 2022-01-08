package memorystorage

import (
	"testing"

	"github.com/LebedevNd/go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	memstorage := New()

	err := memstorage.AddEvent(
		"test title",
		"2021-01-01 15:00:00",
		"2021-01-01 17:00:00",
		"test description",
		42,
		500,
	)
	require.Nil(t, err)
	require.Len(t, memstorage.events, 1)

	err = memstorage.AddEvent(
		"test title",
		"2021-01-01 16:00:00",
		"2021-01-01 18:00:00",
		"test description",
		42,
		500,
	)

	require.Error(t, err)
	require.ErrorIs(t, err, storage.BusyDateError{
		Title:    "test title",
		DateFrom: "2021-01-01 15:00:00",
		DateTo:   "2021-01-01 17:00:00",
	},
	)

	events, err := memstorage.GetEvents()
	require.Nil(t, err)
	require.Len(t, events, 1)

	err = memstorage.UpdateEvent(
		events[0].ID,
		"updated",
		events[0].DateTimeStart,
		events[0].DateTimeEnd,
		events[0].Description,
		events[0].UserID,
		events[0].Delay,
	)
	require.Nil(t, err)
	require.Equal(t, memstorage.events[0].Title, "updated")

	err = memstorage.DeleteEvent(events[0].ID)
	require.Nil(t, err)

	events, err = memstorage.GetEvents()
	require.Nil(t, err)
	require.Len(t, events, 0)
}

func TestStorageListing(t *testing.T) {
	memstorage := New()

	err := memstorage.AddEvent(
		"test title",
		"2021-01-01 15:00:00",
		"2021-01-01 17:00:00",
		"test description",
		42,
		500,
	)
	require.Nil(t, err)

	err = memstorage.AddEvent(
		"test title",
		"2021-01-02 15:00:00",
		"2021-01-02 17:00:00",
		"test description",
		42,
		500,
	)
	require.Nil(t, err)

	err = memstorage.AddEvent(
		"test title",
		"2021-01-22 15:00:00",
		"2021-01-22 17:00:00",
		"test description",
		42,
		500,
	)
	require.Nil(t, err)

	events, _ := memstorage.GetEventsByDay("2021-01-01")
	require.Len(t, events, 1)

	events, _ = memstorage.GetEventsByWeek("2021-01-01")
	require.Len(t, events, 2)

	events, _ = memstorage.GetEventsByMonth("2021-01-01")
	require.Len(t, events, 3)
}
