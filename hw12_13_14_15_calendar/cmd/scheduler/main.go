package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/rabbitmq"
	memorystorage "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage/memory"
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

	storage := memorystorage.New()
	producer, err := rabbitmq.NewProducer(config.MQ)
	if err != nil {
		logg.Error("cannot init producer" + err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler := NewScheduler(*logg, storage, *producer, config.MQ.Interval)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		signal.Stop(signals)

		if err = producer.CloseChannel(); err != nil {
			logg.Error(err.Error())
		}

		if err = producer.CloseConn(); err != nil {
			logg.Error(err.Error())
		}
	}()

	scheduler.Run(ctx, logg, storage, *producer, config.MQ.Interval)
}
