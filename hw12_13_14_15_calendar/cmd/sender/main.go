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
	"github.com/streadway/amqp"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/sender/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := NewSender(configFile)
	if err != nil {
		log.Fatalf("can't get config: %v", err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatalf("can't start logger: %v", err)
	}

	ctx := context.Background()

	consumer, err := rabbitmq.NewConsumer(ctx, config.MQ, logg)
	if err != nil {
		logg.Error("cannot NewConsumer: " + err.Error())
		return
	}

	var delivery <-chan amqp.Delivery
	if delivery, err = consumer.OpenChannel(); err != nil {
		logg.Error("cannot OpenChannel: " + err.Error())
		return
	}
	logg.Info("channel opened")

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		signal.Stop(signals)

		if err = consumer.Close(); err != nil {
			logg.Error(err.Error())
			return
		}
	}()

	done := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		logg.Info("consumer read msg")
		consumer.ReadMsg(delivery, done)
	}()
	wg.Wait()
}
