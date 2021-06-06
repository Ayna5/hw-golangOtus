package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/rabbitmq"
	sqlstorage "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/sheduler/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := NewSheduler(configFile)
	if err != nil {
		log.Fatalf("can't get config1: %v", err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatalf("can't start logger: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage, err := sqlstorage.New(ctx, config.DB.User, config.DB.Password, config.DB.Host, config.DB.Name, config.DB.Port)
	if err != nil {
		logg.Error(err.Error())
		return
	}
	producer, err := rabbitmq.NewProducer(config.MQ)
	if err != nil {
		logg.Error("cannot init producer" + err.Error())
		return
	}

	scheduler := NewScheduler(ctx, *logg, storage, *producer, config.MQ.Interval)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		signal.Stop(signals)

		if err = producer.Close(); err != nil {
			logg.Error(err.Error())
			return
		}
	}()

	done := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		logg.Info("scheduler is running...")
		scheduler.Run(ctx, logg, storage, *producer, config.MQ.Interval, done)
	}()
	wg.Wait()
}
