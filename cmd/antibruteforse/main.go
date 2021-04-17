package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gadzira/anti-bruteforce/internal/app"
	"github.com/gadzira/anti-bruteforce/internal/logger"
	internalhttp "github.com/gadzira/anti-bruteforce/internal/server/http"
	"github.com/gadzira/anti-bruteforce/internal/storage"
	"go.uber.org/zap"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()
	fmt.Println("configFile:", configFile)
	config := NewConfig(configFile)
	l := logger.New(
		config.Logger.LogFile,
		config.Logger.Level,
		config.Logger.MaxSize,
		config.Logger.MaxBackups,
		config.Logger.MaxAge,
		config.Logger.LocalTime,
		config.Logger.Compress,
	)
	logg := l.InitLogger()
	adr := fmt.Sprintf(":%s", config.Server.Port)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bs := storage.New(config.Storage.N, config.Storage.M, config.Storage.K, config.Storage.TTL, logg)

	a := app.New(logg, &bs)
	server := internalhttp.NewServer(logg, a)
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP)

		select {
		case <-ctx.Done():
			return
		case <-signals:
		}

		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	if err := server.Start(ctx, adr); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	logg.Info("server is riseup on %s ...", zap.String("port", adr))
}
