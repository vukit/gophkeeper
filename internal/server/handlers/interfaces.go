package handlers

import (
	"context"
	"io"

	"github.com/vukit/gophkeeper/internal/server/model"
)

// RepoDB интерфейс репозитория базы данных, который используется сервером в обработчиках
type RepoDB interface {
	SaveUser(ctx context.Context, user model.User) (id int, err error)
	FindUser(ctx context.Context, user model.User) (id int, err error)

	SaveLogin(ctx context.Context, login *model.Login) (err error)
	DeleteLogin(ctx context.Context, login *model.Login) (err error)
	FindLogins(ctx context.Context, login model.User) (logins []model.Login, err error)

	SaveCard(ctx context.Context, card *model.Card) (err error)
	DeleteCard(ctx context.Context, card *model.Card) (err error)
	FindCards(ctx context.Context, user model.User) (cards []model.Card, err error)

	SaveFile(ctx context.Context, file *model.File) (err error)
	DeleteFile(ctx context.Context, file *model.File) (err error)
	FindFiles(ctx context.Context, user model.User) (files []model.File, err error)
	FindFile(ctx context.Context, fileID, userID int) (file *model.File, err error)

	Close() error
}

// RepoFile интерфейс репозитория файлов, который используется сервером в обработчиках
type RepoFile interface {
	SaveFile(ctx context.Context, src io.ReadCloser) (filePath string, err error)
	DeleteFile(ctx context.Context, filePath string) (err error)
	GetFile(ctx context.Context, filePath string) (src io.ReadCloser, err error)

	Close() error
}
