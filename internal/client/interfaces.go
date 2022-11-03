package client

import (
	"context"

	"github.com/vukit/gophkeeper/internal/client/model"
	"github.com/vukit/gophkeeper/internal/client/service"
)

// GophKeeperService описывает методы, необходимые для работы всего приложения
type GophKeeperService interface {
	SetCryptoService(cs *service.CryptoService)

	SignUp(context.Context, *model.User) error
	SignIn(context.Context, *model.User) error

	SaveLogin(context.Context, *model.Login) error
	DeleteLogin(context.Context, *model.Login) error
	GetLogins(context.Context) ([]model.Login, error)

	SaveCard(context.Context, *model.Card) error
	DeleteCard(context.Context, *model.Card) error
	GetCards(context.Context) ([]model.Card, error)

	SaveFile(context.Context, *model.File) error
	DeleteFile(context.Context, *model.File) error
	GetFiles(context.Context) ([]model.File, error)
	DownloadFile(context.Context, *model.File, string) error
}
