package app

import (
	"context"

	"github.com/vukit/gophkeeper/internal/client"
	"github.com/vukit/gophkeeper/internal/client/config"
	"github.com/vukit/gophkeeper/internal/client/logger"
	"github.com/vukit/gophkeeper/internal/client/service"
	"github.com/vukit/gophkeeper/internal/client/tui"
)

func Run(ctx context.Context, cfg *config.Config, mLogger *logger.Logger) {
	var gkService client.GophKeeperService

	switch cfg.ServerProtocol {
	case "http", "https":
		gkService = service.NewHTTPService(cfg, mLogger)
	}

	if cfg.UserInterface == "tui" {
		user, err := tui.Login(ctx, gkService, mLogger)
		if err != nil {
			mLogger.Info(err.Error())

			return
		}

		cryptoService, err := service.NewCryptoService(user)
		if err != nil {
			mLogger.Info(err.Error())

			return
		}

		gkService.SetCryptoService(cryptoService)

		err = tui.Manager(ctx, user, gkService, mLogger, cfg.DownloadFolder)
		if err != nil {
			mLogger.Info(err.Error())

			return
		}
	} else {
		mLogger.Info("Graphical user interface not implemented yet")
	}
}
