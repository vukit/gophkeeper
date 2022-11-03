package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env"
	"github.com/vukit/gophkeeper/internal/client/app"
	"github.com/vukit/gophkeeper/internal/client/config"
	"github.com/vukit/gophkeeper/internal/client/logger"
)

func main() {
	cfg := config.Config{}
	flag.StringVar(&cfg.ServerAddress, "s", "localhost:8080", "server address")
	flag.StringVar(&cfg.ServerProtocol, "p", "http", "server protocol (http|https)")
	flag.StringVar(&cfg.LogFile, "l", "client.log", "logging file")
	flag.StringVar(&cfg.UserInterface, "u", "tui", "user interface (tui|gui)")
	flag.StringVar(&cfg.DownloadFolder, "d", "gophkeeper", "folder for downloaded files")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)

		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	logFile, err := os.Create(cfg.LogFile)
	if err != nil {
		log.Println(err)

		return
	}
	defer logFile.Close()

	mLogger := logger.NewLogger(logFile)

	if _, err := os.Stat(cfg.DownloadFolder); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(cfg.DownloadFolder, os.ModePerm)
		if err != nil {
			log.Println(err)

			return
		}
	}

	app.Run(ctx, &cfg, mLogger)
}
