package tui

import (
	"context"
	"time"

	"github.com/rivo/tview"
	"github.com/vukit/gophkeeper/internal/client"
	"github.com/vukit/gophkeeper/internal/client/model"
)

// Logins компонент реализует текстовый интерфейс CRUD для пар username/password
func (r *TUI) Logins(ctx context.Context, user *model.User, service client.GophKeeperService) *tview.Flex {
	login := &model.Login{}

	layout := tview.NewFlex()

	list := tview.NewList()
	list.SetTitle("[ Logins ]").SetBorder(true).SetBorderPadding(1, 1, 1, 1)
	list.ShowSecondaryText(false)

	form := tview.NewForm()
	login.ID = 0
	login.Username = ""
	login.Password = ""
	login.MetaInfo = ""
	setupLoginForm(ctx, form, login, r, service)

	layout.AddItem(list, 35, 0, false).AddItem(form, 0, 1, true)

	go loginsUpdateList(ctx, login, list, form, r, service)

	return layout
}

func setupLoginForm(
	ctx context.Context,
	form *tview.Form,
	login *model.Login,
	r *TUI,
	service client.GophKeeperService,
) {
	form.Clear(true)
	form.
		AddInputField("Username", login.Username, 30, nil, func(text string) { login.Username = text }).
		AddPasswordField("Password", login.Password, 30, '*', func(text string) { login.Password = text }).
		AddInputField("Metainfo", login.MetaInfo, 30, nil, func(text string) { login.MetaInfo = text }).
		AddButton("Save", func() {
			err := login.Validate()
			if err != nil {
				r.alertChannel <- err.Error()

				return
			}

			err = service.SaveLogin(ctx, login)
			if err != nil {
				r.alertChannel <- err.Error()

				return
			}

			r.alertChannel <- "login was successfully saved"
			login.ID = 0
			login.Username = ""
			login.Password = ""
			login.MetaInfo = ""
			setupLoginForm(ctx, form, login, r, service)
		}).
		AddButton("Cancel", func() {
			login.ID = 0
			login.Username = ""
			login.Password = ""
			login.MetaInfo = ""
			setupLoginForm(ctx, form, login, r, service)
		})

	if idx := form.GetButtonIndex("Delete"); idx != -1 {
		form.RemoveButton(idx)
	}

	form.SetBorder(true).SetTitle("[ New login ]").SetTitleAlign(tview.AlignLeft).SetBorderPadding(1, 0, 2, 0)

	r.app.SetFocus(form)
}

func loginsUpdateList(
	ctx context.Context,
	login *model.Login,
	list *tview.List,
	form *tview.Form,
	r *TUI,
	service client.GophKeeperService,
) {
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			logins, err := service.GetLogins(ctx)
			if err != nil {
				r.alertChannel <- err.Error()

				continue
			}

			currentItemIndex := list.GetCurrentItem()

			list.Clear()

			for _, currentLogin := range logins {
				currentLogin := currentLogin
				list.AddItem(currentLogin.MetaInfo, "", rune(0), func() {
					login = &currentLogin
					setupLoginForm(ctx, form, login, r, service)
					form.SetTitle("[ Edit login ]")
					if idx := form.GetButtonIndex("Delete"); idx == -1 {
						form.AddButton("Delete", func() {
							errService := service.DeleteLogin(ctx, login)
							if errService != nil {
								r.alertChannel <- errService.Error()

								return
							}
							r.alertChannel <- "login was successfully deleted"
							login.ID = 0
							login.Username = ""
							login.Password = ""
							login.MetaInfo = ""
							setupLoginForm(ctx, form, login, r, service)
						})
					}
				})
			}

			list.SetCurrentItem(currentItemIndex)

			r.app.Draw()
		case <-ctx.Done():
			ticker.Stop()

			return
		}
	}
}
