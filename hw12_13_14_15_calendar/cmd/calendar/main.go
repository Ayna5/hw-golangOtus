package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/app"
	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/logger"
	grpc2 "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage/sql"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var calendar *app.App
	if config.DB.Mem {
		storage := memorystorage.New()
		calendar = app.New(logg, storage)
	} else {
		storage, err := sqlstorage.New(ctx, config.DB.User, config.DB.Password, config.DB.Host, config.DB.Name, config.DB.Port)
		if err != nil {
			logg.Error(err.Error())
		}
		calendar = app.New(logg, storage)
	}

	grpc, err := grpc2.NewServer(logg, calendar, config.Server.Grpc)
	if err != nil {
		log.Fatal("cannot init grpc server") //nolint:gocritic
	}
	defer grpc.Stop() //nolint:errcheck

	server := internalhttp.NewServer(config.Server.HTTP, *logg, calendar)

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
			logg.Error("failed to stop query.http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err = grpc.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
		}
	}()

	go func() {
		defer wg.Done()
		if err = server.Start(ctx); err != nil {
			logg.Error("failed to start query.http server: " + err.Error())
		}
	}()
	wg.Wait()
}
