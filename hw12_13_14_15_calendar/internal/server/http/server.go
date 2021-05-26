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
	api     *app.App
}

func NewServer(host string, log logger.Logger, api *app.App) *Server {
	return &Server{
		Address: host,
		logger:  log,
		api:     api,
	}
}

func (s *Server) Start(ctx context.Context) error {
	router := mux.NewRouter()
	router.HandleFunc("/hello", s.helloWorld).Methods("GET")
	router.HandleFunc("/create", s.createEvent).Methods("POST")
	router.HandleFunc("/update", s.updateEvent).Methods("POST")
	router.HandleFunc("/get", s.getEvents).Methods("GET")
	router.HandleFunc("/delete", s.deleteEvent).Methods("POST")
	router.Use(s.loggingMiddleware)

	s.server = &http.Server{
		Addr:         s.Address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "start server error")
	}

	<-ctx.Done()
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

func (s *Server) helloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	var event storage.Event
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("cannot decode to struct" + err.Error()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if err = s.api.CreateEvent(event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("cannot create event" + err.Error()))
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

func (s *Server) updateEvent(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	var event storage.Event
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("cannot decode to struct" + err.Error()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if err = s.api.UpdateEvent(event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("cannot update event" + err.Error()))
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

func (s *Server) getEvents(w http.ResponseWriter, r *http.Request) {
	var event storage.Event
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("cannot decode to struct" + err.Error()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	events, err := s.api.GetEvents(event.StartData, event.EndData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("cannot get events" + err.Error()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if len(events) == 0 {
		w.WriteHeader(http.StatusNotFound)
		_, err = w.Write([]byte("events not found" + err.Error()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	result, err := json.Marshal(events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("cannot execute marshal" + err.Error()))
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

func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	var event storage.Event
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("cannot decode to struct" + err.Error()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if err = s.api.DeleteEvent(event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("cannot delete event" + err.Error()))
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
