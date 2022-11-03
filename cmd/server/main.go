package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"github.com/vukit/gophkeeper/internal/server/config"
	"github.com/vukit/gophkeeper/internal/server/logger"
	"github.com/vukit/gophkeeper/internal/server/repositories/localfiles"
	"github.com/vukit/gophkeeper/internal/server/repositories/postgresql"
	"github.com/vukit/gophkeeper/internal/server/router"
	"github.com/vukit/gophkeeper/internal/server/utils"
	"golang.org/x/sync/errgroup"
)

// @Title GophKeeper API
// @Version 1.0
// @Contact.name Mark Vaisman
// @License.name MIT
// @License.url https://github.com/vukit/gophkeeper/blob/main/LICENSE
// @BasePath /api
func main() {
	mLogger := logger.NewLogger(os.Stderr)

	mConfig := config.Config{}
	flag.StringVar(&mConfig.Address, "a", "localhost:8080", "server address")
	flag.StringVar(&mConfig.Protocol, "p", "http", "server protocol (http|https)")
	flag.StringVar(&mConfig.TLSCertificate, "tlsc", "", "path to TLS certificate")
	flag.StringVar(&mConfig.TLSPrivateKey, "tlspk", "", "path to TLS private key")
	flag.StringVar(&mConfig.DataBaseURI, "d", "postgres://postgres:postgres@localhost:5432/gophkeeper?sslmode=disable", "database uri")
	flag.StringVar(&mConfig.FileStorage, "s", "storage", "file storage path")
	flag.Parse()

	err := env.Parse(&mConfig)
	if err != nil {
		mLogger.Fatal(err.Error())
	}

	utils.MigrationUp("file://internal/server/migrations/", mConfig.DataBaseURI, mLogger)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	mRepoDB, err := postgresql.NewRepo(mConfig.DataBaseURI)
	if err != nil {
		mLogger.Fatal(err.Error())
	}
	defer mRepoDB.Close()

	mRepoFile, err := localfiles.NewRepo(mConfig.FileStorage)
	if err != nil {
		mLogger.Fatal(err.Error())
	}
	defer mRepoFile.Close()

	errGroup, errGroupCtx := errgroup.WithContext(ctx)

	switch mConfig.Protocol {
	case "http", "https":
		mRouter, err := router.NewRouter(ctx, mRepoDB, mRepoFile, mLogger)
		if err != nil {
			mLogger.Fatal(err.Error())
		}

		mServer := &http.Server{Addr: mConfig.Address, Handler: mRouter, ReadHeaderTimeout: time.Second}

		errGroup.Go(func() error {
			switch mConfig.Protocol {
			case "http":
				return mServer.ListenAndServe()
			case "https":
				return mServer.ListenAndServeTLS(mConfig.TLSCertificate, mConfig.TLSPrivateKey)
			default:
				return nil
			}
		})

		errGroup.Go(func() error {
			<-errGroupCtx.Done()

			return mServer.Shutdown(context.Background())
		})
	default:
		mLogger.Fatal(fmt.Sprintf("unknown protocol: %s", mConfig.Protocol))
	}

	mLogger.Info("Server started")

	if err := errGroup.Wait(); err != nil {
		mLogger.Info(err.Error())
	}
}
