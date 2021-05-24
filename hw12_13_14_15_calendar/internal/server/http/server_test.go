package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/app"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/logger"
	storage2 "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

var (
	ctx = context.Background()
)

func TestServerGRPC(t *testing.T) {
	var l logger.Logger
	host := "localhost:8080"
	storage := memorystorage.New()
	api := app.New(l, storage)
	server := NewServer(host, l, api)
	defer server.Stop(ctx)

	t.Run("Create event", func(t *testing.T) {
		newEvent := storage2.Event{
			ID:          "2",
			Title:       "event2",
			StartData:   time.Now(),
			EndData:     time.Now().AddDate(0, 0, 2),
			Description: "event2",
			OwnerID:     "2",
			RemindIn:    "3",
		}
		data, err := json.Marshal(&newEvent)
		require.NoError(t, err)

		resp, err := http.Post("http://"+host+"/create", "application/json", bytes.NewReader(data))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
	t.Run("Update event", func(t *testing.T) {
		newEvent := storage2.Event{
			ID:          "2",
			Title:       "event3",
			StartData:   time.Now(),
			EndData:     time.Now().AddDate(0, 0, 2),
			Description: "event3",
			OwnerID:     "3",
			RemindIn:    "3",
		}
		data, err := json.Marshal(&newEvent)
		require.NoError(t, err)

		resp, err := http.Post("http://"+host+"/update", "application/json", bytes.NewReader(data))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
	t.Run("Get events", func(t *testing.T) {
		newEvent := storage2.Event{
			StartData: time.Now().Add(-1 * time.Hour),
			EndData:   time.Now().AddDate(0, 0, 3),
		}
		data, err := json.Marshal(&newEvent)
		require.NoError(t, err)

		_, err = http.NewRequest(http.MethodGet, "http://"+host+"/get", bytes.NewReader(data))
		require.NoError(t, err)
	})
	t.Run("Delete event", func(t *testing.T) {
		newEvent := storage2.Event{
			ID:          "2",
			Title:       "event3",
			StartData:   time.Now(),
			EndData:     time.Now().AddDate(0, 0, 2),
			Description: "event3",
			OwnerID:     "3",
			RemindIn:    "3",
		}
		data, err := json.Marshal(&newEvent)
		require.NoError(t, err)

		resp, err := http.Post("http://"+host+"/delete", "application/json", bytes.NewReader(data))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
