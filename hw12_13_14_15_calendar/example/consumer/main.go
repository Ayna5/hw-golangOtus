package main

import (
	"context"
	"flag"
	"log"

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

	ctx := context.Background()

	config, err := NewSender(configFile)
	if err != nil {
		log.Fatalf("can't get config1: %v", err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatalf("can't start logger: %v", err)
	}

	consumer, err := rabbitmq.NewConsumer(ctx, config.MQ, logg)
	if err != nil {
		logg.Error("cannot NewConsumer: " + err.Error())
	}
	defer consumer.Close()

	logg.Info("open channel")
	var delivery <-chan amqp.Delivery
	if delivery, err = consumer.OpenChannel(); err != nil {
		logg.Error("cannot OpenChannel: " + err.Error())
	}

	logg.Info("consumer read msg")
	consumer.ReadMsg(delivery, make(chan error))
}
