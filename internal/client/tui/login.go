package tui

import (
	"context"
	"errors"

	"github.com/rivo/tview"
	"github.com/vukit/gophkeeper/internal/client"
	"github.com/vukit/gophkeeper/internal/client/logger"
	"github.com/vukit/gophkeeper/internal/client/model"
)

// Login предоставляет форму аутентификации и регистрации пользователя,
// возвращает идентифицированного пользователя либо ошибку
func Login(ctx context.Context, service client.GophKeeperService, mLogger *logger.Logger) (user *model.User, err error) {
	user = &model.User{}
	tvApp := tview.NewApplication()

	alert, _ := Alert(ctx, 1, 0, 0, 0, tview.AlignCenter, 0)

	form := tview.NewForm().
		AddInputField("Username", user.Username, 23, nil, func(text string) { user.Username = text }).
		AddPasswordField("Password", user.Password, 23, '*', func(text string) { user.Password = text }).
		AddButton("Sign in", func() {
			err = user.Validate()
			if err != nil {
				alert.SetText(err.Error())

				return
			}

			err = service.SignIn(ctx, user)
			if err != nil {
				alert.SetText(err.Error())

				return
			}

			tvApp.Stop()
		}).
		AddButton("Sign up", func() {
			err = user.Validate()
			if err != nil {
				alert.SetText(err.Error())

				return
			}

			err = service.SignUp(ctx, user)
			if err != nil {
				alert.SetText(err.Error())

				return
			}

			tvApp.Stop()
		}).
		AddButton("Quit", func() {
			err = errors.New("canceled authentication")
			tvApp.Stop()
		})
	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(alert, 2, 0, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(form, 8, 1, true).
				AddItem(Copyright(tview.AlignCenter), 1, 0, false),
				36, 0, true).
			AddItem(tview.NewBox(), 0, 1, false),
			0, 1, true)

	if errTVApp := tvApp.SetRoot(layout, true).EnableMouse(true).Run(); errTVApp != nil {
		return nil, errTVApp
	}

	return user, err
}
