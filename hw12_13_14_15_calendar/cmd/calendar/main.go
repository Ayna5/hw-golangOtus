package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/app"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/logger"
	grpc2 "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("can't get config: %v", err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatalf("can't start logger: %v", err)
	}

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	grpc, err := grpc2.NewServer(logg, calendar, config.Server.Grpc)
	if err != nil {
		log.Fatal("cannot init grpc server")
	}
	defer grpc.Stop() //nolint:errcheck

	server := internalhttp.NewServer(config.Server.HTTP, *logg, calendar)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		signal.Stop(signals)
		cancel()

		ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err = grpc.Stop(); err != nil {
			logg.Error("cannot close connection: " + err.Error())
		}
		if err = server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")
	if err = grpc.Start(); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
		defer cancel()
	}

	if err = server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		defer cancel()
	}
}
