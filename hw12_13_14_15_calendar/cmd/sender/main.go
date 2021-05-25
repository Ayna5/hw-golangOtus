package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/rabbitmq"
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

	consumer, err := rabbitmq.NewConsumer(config.MQ)
	if err != nil {
		logg.Error("cannot NewConsumer: " + err.Error())
	}

	if err = consumer.OpenChannel(); err != nil {
		logg.Error("cannot OpenChannel: " + err.Error())
	}

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		signal.Stop(signals)

		if err = consumer.CloseChannel(); err != nil {
			logg.Error(err.Error())
		}

		if err = consumer.CloseConn(); err != nil {
			logg.Error(err.Error())
		}
	}()
}
