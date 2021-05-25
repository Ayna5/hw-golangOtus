package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/app"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	ctx = context.Background()
)

func newTestServer(t *testing.T) *httptest.Server {
	var logg logrus.Level
	log, err := logger.New(logg.String(), "logrus.log")
	require.NoError(t, err)

	storage := memorystorage.New()
	a := app.New(log, storage)
	api := NewServer(":5555", *log, a)
	router := api.Route()
	require.NoError(t, err)

	return httptest.NewServer(router)
}

func TestServerHTTP(t *testing.T) {
	svc := newTestServer(t)

	t.Run("Create event", func(t *testing.T) {
		newEvent := storage.Event{
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

		resp, err := http.Post(svc.URL+"/create", "application/json", bytes.NewReader(data))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
	t.Run("Update event", func(t *testing.T) {
		newEvent := storage.Event{
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

		resp, err := http.Post(svc.URL+"/update", "application/json", bytes.NewReader(data))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
	t.Run("Get events", func(t *testing.T) {
		newEvent := storage.Event{
			StartData: time.Now().Add(-1 * time.Hour),
			EndData:   time.Now().AddDate(0, 0, 3),
		}
		data, err := json.Marshal(&newEvent)
		require.NoError(t, err)

		_, err = http.NewRequest(http.MethodGet, svc.URL+"/get", bytes.NewReader(data))
		require.NoError(t, err)
	})
	t.Run("Delete event", func(t *testing.T) {
		newEvent := storage.Event{
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

		resp, err := http.Post(svc.URL+"/delete", "application/json", bytes.NewReader(data))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
