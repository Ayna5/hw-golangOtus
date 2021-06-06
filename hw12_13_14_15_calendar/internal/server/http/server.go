package internalhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/app"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Server struct {
	Address string
	server  *http.Server
	logger  logger.Logger
}

type Handlers struct {
	logger logger.Logger
	api    *app.App
}

func NewServer(host string, log logger.Logger, api *app.App) *Server {
	mux := route(log, api)
	server := &http.Server{
		Addr:         host,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		Address: host,
		server:  server,
		logger:  log,
	}
}

func route(l logger.Logger, app *app.App) *mux.Router {
	handler := &Handlers{
		logger: l,
		api:    app,
	}
	router := mux.NewRouter()
	router.HandleFunc("/hello", handler.helloWorld).Methods("GET")
	router.HandleFunc("/create", handler.createEvent).Methods("POST")
	router.HandleFunc("/update", handler.updateEvent).Methods("POST")
	router.HandleFunc("/get", handler.getEvents).Methods("GET")
	router.HandleFunc("/delete", handler.deleteEvent).Methods("POST")
	router.Use(handler.loggingMiddleware)

	return router
}

func (s *Server) Start(ctx context.Context) error {
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "start server error")
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return errors.New("server is nil")
	}
	if err := s.server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "stop server error")
	}
	return nil
}

func (h *Handlers) helloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}

func (h *Handlers) createEvent(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	var event storage.Event
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("cannot decode to struct: " + err.Error()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if err = h.api.CreateEvent(event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Событие создано"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) updateEvent(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	var event storage.Event
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("cannot decode to struct: " + err.Error()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if err = h.api.UpdateEvent(event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Событие обновлено"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) getEvents(w http.ResponseWriter, r *http.Request) {
	var event storage.Event
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("cannot decode to struct: " + err.Error()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	events, err := h.api.GetEvents(r.Context(), event.StartData, event.EndData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if events == nil && err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, err = w.Write([]byte("events not found: " + err.Error()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if events == nil {
		events = []storage.Event{}
	}
	result, err := json.Marshal(events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("cannot execute marshal: " + err.Error()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) deleteEvent(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	var event storage.Event
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("cannot decode to struct: " + err.Error()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if err = h.api.DeleteEvent(event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Событие удалено"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
