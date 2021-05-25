package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/configs"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/rabbitmq"
	memorystorage "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage/memory"
)

type Scheduler struct {
	Log      logger.Logger
	Storage  *memorystorage.Storage
	Producer rabbitmq.Producer
	Interval time.Duration
}

func NewScheduler(logger logger.Logger, storage *memorystorage.Storage, producer rabbitmq.Producer, interval time.Duration) *Scheduler {
	return &Scheduler{Log: logger, Storage: storage, Producer: producer, Interval: interval}
}

func startWorker(ctx context.Context, interval time.Duration, f func()) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			f()
		case <-ctx.Done():
			return
		}
	}
}

func (s *Scheduler) Run(ctx context.Context, log *logger.Logger, sql *memorystorage.Storage, producer rabbitmq.Producer, t time.Duration) {
	go startWorker(ctx, s.Interval, func() {
		s.makeTasks(log, sql, producer, t)
	})
}

func (s *Scheduler) makeTasks(log *logger.Logger, sql *memorystorage.Storage, producer rabbitmq.Producer, t time.Duration) {
	from := time.Now()
	to := from.Add(t)
	oneYearLater := from.AddDate(-1, 0, 0)

	_ = producer.OpenChannel()
	events, err := sql.GetEvents(from, to)
	if err != nil {
		log.Error(err.Error())
	}

	for _, e := range events {
		data, err := json.Marshal(configs.ConvertToMQNotification(*e))
		if err != nil {
			log.Error(err.Error())
			continue
		}
		err = s.Producer.Send(data)
		if err != nil {
			log.Error(err.Error())
		}

		if e.StartData.Before(oneYearLater) {
			err = sql.DeleteEvent(*e)
			if err != nil {
				log.Error(err.Error())
			}
		}
	}
}
