package tui

import (
	"context"

	"github.com/rivo/tview"
	"github.com/vukit/gophkeeper/internal/client"
	"github.com/vukit/gophkeeper/internal/client/logger"
	"github.com/vukit/gophkeeper/internal/client/model"
)

// TUI структура данных приложения текстового интерфейса
type TUI struct {
	app          *tview.Application
	alertChannel chan string
	mLogger      *logger.Logger
}

// Component структура данных компонента текстового интерфейса
type Component struct {
	name   string
	layout *tview.Flex
}

// Manager реализует механизм переключения компонентов текстового интерфейса
func Manager(
	ctx context.Context,
	user *model.User,
	service client.GophKeeperService,
	mLogger *logger.Logger,
	downloadFolder string,
) (err error) {
	alertBox, alertChannel := Alert(ctx, 0, 0, 0, 0, tview.AlignLeft, 3)

	tui := TUI{app: tview.NewApplication(), alertChannel: alertChannel, mLogger: mLogger}

	components := []Component{
		{"Logins", tui.Logins(ctx, user, service)},
		{"Cards", tui.Cards(ctx, user, service)},
		{"Files", tui.Files(ctx, user, service, downloadFolder)},
	}

	pages := tview.NewPages()

	for _, component := range components {
		func(component Component) {
			pages.AddPage(component.name,
				component.layout,
				true,
				component.name == components[0].name)
		}(component)
	}

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewGrid().SetColumns(-1, 35).
			AddItem(alertBox, 0, 0, 1, 1, 0, 0, false).
			AddItem(Copyright(tview.AlignRight), 0, 1, 1, 1, 0, 0, false),
			1, 0, false).
		AddItem(pages, 0, 1, true).
		AddItem(tview.NewGrid().SetColumns(10, 10, 10, 10, -1, 10).SetGap(0, 1).
			AddItem(tview.NewButton(components[0].name).SetSelectedFunc(func() {
				pages.SwitchToPage(components[0].name)
				tui.app.SetFocus(components[0].layout)
			}), 0, 0, 1, 1, 0, 0, false).
			AddItem(tview.NewButton(components[1].name).SetSelectedFunc(func() {
				pages.SwitchToPage(components[1].name)
				tui.app.SetFocus(components[1].layout)
			}), 0, 1, 1, 1, 0, 0, false).
			AddItem(tview.NewButton(components[2].name).SetSelectedFunc(func() {
				pages.SwitchToPage(components[2].name)
				tui.app.SetFocus(components[2].layout)
			}), 0, 2, 1, 1, 0, 0, false).
			AddItem(tview.NewBox(), 0, 4, 1, 1, 0, 0, false).
			AddItem(tview.NewButton("Quit").SetSelectedFunc(func() {
				tui.app.Stop()
			}), 0, 5, 1, 1, 0, 0, false),
			1, 0, false)

	if errTVApp := tui.app.SetRoot(layout, true).EnableMouse(true).Run(); errTVApp != nil {
		return errTVApp
	}

	return nil
}
