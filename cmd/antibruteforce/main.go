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
	conf "github.com/gadzira/anti-bruteforce/internal/config"
	"github.com/gadzira/anti-bruteforce/internal/database"
	"github.com/gadzira/anti-bruteforce/internal/helpers"
	"github.com/gadzira/anti-bruteforce/internal/logger"
	internalhttp "github.com/gadzira/anti-bruteforce/internal/server/http"
	"github.com/gadzira/anti-bruteforce/internal/storage"
	"go.uber.org/zap"
)

var (
	configFile  string
	resetBucket string
	addWhite    string
	addBlack    string
	delWhite    string
	delBlack    string
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
	flag.StringVar(&resetBucket, "reset-bucket", "", "Use: ./abf reset-bucket=ip/login for reset bucket")
	flag.StringVar(&addWhite, "add-white", "", "Use: ./abf add-white=ip:mask for add IP to white list")
	flag.StringVar(&addBlack, "add-black", "", "Use: ./abf add-black=ip:mask for add IP to black list")
	flag.StringVar(&delWhite, "del-white", "", "Use: ./abf del-white=ip:mask for delete IP from white list")
	flag.StringVar(&delBlack, "del-black", "", "Use: ./abf del-black=ip:mask for delete IP from black list")
}

func main() {
	flag.Parse()
	config := conf.NewFromFile(configFile)
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
	adr := fmt.Sprintf("127.0.0.1:%s", config.Server.Port)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sql := database.New(logg)
	err := sql.Connect(ctx, config.DataBase.DSN)
	if err != nil {
		logg.Fatal("can't connect to DB: %s\n", zap.String("err", err.Error()))
	}

	bs := storage.New(config.Storage.N, config.Storage.M, config.Storage.K, config.Storage.TTL, logg)
	a := app.New(ctx, logg, sql, &bs)
	server := internalhttp.NewServer(logg, a)

	switch {
	case resetBucket != "":
		/*
			That will not work, because when the application is running,
			memory storage is empty.
		*/
		bs.ResetBucket(resetBucket)
	case addWhite != "":
		e := helpers.MakeEntry(addWhite, "white")
		err := sql.AddToList(ctx, e)
		if err != nil {
			a.Logger.Fatal("can't add to list\n" + err.Error())
		}
	case addBlack != "":
		e := helpers.MakeEntry(addBlack, "black")
		err := sql.AddToList(ctx, e)
		if err != nil {
			a.Logger.Fatal("can't add to list\n" + err.Error())
		}
	case delWhite != "":
		e := helpers.MakeEntry(addBlack, "black")
		err := sql.RemoveFromList(ctx, e)
		if err != nil {
			a.Logger.Fatal("can't remove from list\n" + err.Error())
		}
	case delBlack != "":
		e := helpers.MakeEntry(addBlack, "black")
		err := sql.RemoveFromList(ctx, e)
		if err != nil {
			a.Logger.Fatal("can't remove from list\n" + err.Error())
		}
	default:
		logg.Info("additional arguments not given yet")
	}

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
