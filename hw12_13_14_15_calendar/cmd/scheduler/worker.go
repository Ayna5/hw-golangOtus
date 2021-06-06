package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/configs"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/rabbitmq"
	sqlstorage "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage/sql"
)

type Scheduler struct {
	Ctx      context.Context
	Log      logger.Logger
	Storage  *sqlstorage.Storage
	Producer rabbitmq.Producer
	Interval time.Duration
}

func NewScheduler(ctx context.Context, logger logger.Logger, storage *sqlstorage.Storage, producer rabbitmq.Producer, interval time.Duration) *Scheduler {
	return &Scheduler{Ctx: ctx, Log: logger, Storage: storage, Producer: producer, Interval: interval}
}

func startWorker(ctx context.Context, done chan error, interval time.Duration, f func()) {
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ctx done")
				close(done)
			case <-ticker.C:
				fmt.Println("ticker")
				f()
			}
		}
	}()
}

func (s *Scheduler) Run(ctx context.Context, log *logger.Logger, sql *sqlstorage.Storage, producer rabbitmq.Producer, t time.Duration, done chan error) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go startWorker(ctx, done, s.Interval, func() {
		defer wg.Done()
		s.makeTasks(ctx, log, sql, producer, t)
	})
	wg.Wait()
}

func (s *Scheduler) makeTasks(ctx context.Context, log *logger.Logger, sql *sqlstorage.Storage, producer rabbitmq.Producer, t time.Duration) {
	from := time.Now().Add((-t) * time.Minute).UTC()
	to := time.Now().UTC()
	oneYearLater := from.AddDate(-1, 0, 0)

	err := s.Producer.OpenChannel()
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer producer.Close()

	events, err := sql.GetEvents(ctx, from, to)
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, e := range events {
		data, err := json.Marshal(configs.ConvertToMQNotification(e))
		if err != nil {
			log.Error(err.Error())
			continue
		}
		log.Info(string(data))
		err = s.Producer.Send(data)
		if err != nil {
			log.Error(err.Error())
			return
		}

		if e.StartData.Before(oneYearLater) {
			err = sql.DeleteEvent(e)
			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	}
}
