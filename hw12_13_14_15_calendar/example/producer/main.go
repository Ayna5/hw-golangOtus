package main

import (
	"flag"
	"log"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/rabbitmq"
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

	producer, err := rabbitmq.NewProducer(config.MQ)
	if err != nil {
		logg.Error("cannot init producer" + err.Error())
	}

	if err := producer.OpenChannel(); err != nil {
		logg.Error(err.Error())
	}
	defer producer.Close()

	if err := producer.Send([]byte(
		// "{\"id\": \"4\",\"title\": \"eventNew\",\"startData\": \"2021-05-30T09:00:00Z\",\"endData\": \"2021-12-30T19:10:25Z\",\"description\": \"event is new\",\"ownerId\": \"12\",\"remindIn\": \"2\"}",
		"{\"id\": \"5\",\"title\": \"eventNew\",\"startData\": \"2021-05-30T09:00:00Z\",\"endData\": \"2021-12-30T19:10:25Z\",\"description\": \"event is new\",\"ownerId\": \"12\",\"remindIn\": \"2\"}",
	)); err != nil {
		logg.Error(err.Error())
	}
}
