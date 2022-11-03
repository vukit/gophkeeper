package tui

import (
	"context"
	"time"

	"github.com/rivo/tview"
	"github.com/vukit/gophkeeper/internal/client"
	"github.com/vukit/gophkeeper/internal/client/model"
)

// Files компонент реализует текстовый интерфейс CRUD для файлов пользователя
func (r *TUI) Files(
	ctx context.Context,
	user *model.User,
	service client.GophKeeperService,
	downloadFolder string,
) *tview.Flex {
	file := &model.File{}

	layout := tview.NewFlex()

	list := tview.NewList()
	list.SetTitle("[ Files ]").SetBorder(true).SetBorderPadding(1, 1, 1, 1)

	form := tview.NewForm()
	file.ID = 0
	file.Path = ""
	file.MetaInfo = ""
	setupFileForm(ctx, form, file, r, service, downloadFolder)

	fileBrowser := r.FileBrowser(form, downloadFolder)

	layout.AddItem(list, 35, 0, false).AddItem(form, 0, 1, true).AddItem(fileBrowser, 60, 0, false)

	go filesUpdateList(ctx, file, list, form, r, service, downloadFolder)

	return layout
}

func setupFileForm(
	ctx context.Context,
	form *tview.Form,
	file *model.File,
	r *TUI,
	service client.GophKeeperService,
	downloadFolder string,
) {
	form.Clear(true)
	form.
		AddInputField("Metainfo", file.MetaInfo, 60, nil, func(text string) { file.MetaInfo = text }).
		AddInputField("File", "", 60, nil, func(text string) { file.Path = text }).
		AddButton("Save", func() {
			err := file.Validate()
			if err != nil {
				r.alertChannel <- err.Error()

				return
			}

			r.alertChannel <- "saving file..."
			err = service.SaveFile(ctx, file)
			if err != nil {
				r.alertChannel <- err.Error()

				return
			}

			r.alertChannel <- "file was successfully saved"
			file.ID = 0
			file.Path = ""
			file.MetaInfo = ""
			setupFileForm(ctx, form, file, r, service, downloadFolder)
		}).
		AddButton("Cancel", func() {
			file.ID = 0
			file.Path = ""
			file.MetaInfo = ""
			setupFileForm(ctx, form, file, r, service, downloadFolder)
		})

	if idx := form.GetButtonIndex("Delete"); idx != -1 {
		form.RemoveButton(idx)
	}

	form.SetBorder(true).SetTitle("[ New file ]").SetTitleAlign(tview.AlignLeft).SetBorderPadding(1, 0, 2, 0)

	r.app.SetFocus(form)
}

func filesUpdateList(
	ctx context.Context,
	file *model.File,
	list *tview.List,
	form *tview.Form,
	r *TUI,
	service client.GophKeeperService,
	downloadFolder string,
) {
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			files, err := service.GetFiles(ctx)
			if err != nil {
				r.alertChannel <- err.Error()

				continue
			}

			currentItemIndex := list.GetCurrentItem()

			list.Clear()

			for _, currentFile := range files {
				currentFile := currentFile
				list.AddItem(currentFile.MetaInfo, currentFile.Name, rune(0), func() {
					file = &currentFile
					setupFileForm(ctx, form, file, r, service, downloadFolder)
					form.SetTitle("[ Edit файл ]")
					if idx := form.GetButtonIndex("Download"); idx == -1 {
						form.AddButton("Download", func() {
							errService := service.DownloadFile(ctx, file, downloadFolder)
							if errService != nil {
								r.alertChannel <- errService.Error()

								return
							}
							r.alertChannel <- "file was successfully downloaded"
							file.ID = 0
							file.Path = ""
							file.MetaInfo = ""
							setupFileForm(ctx, form, file, r, service, downloadFolder)
						})
					}
					if idx := form.GetButtonIndex("Delete"); idx == -1 {
						form.AddButton("Delete", func() {
							errService := service.DeleteFile(ctx, file)
							if errService != nil {
								r.alertChannel <- errService.Error()

								return
							}
							r.alertChannel <- "file was successfully deleted"
							file.ID = 0
							file.Path = ""
							file.MetaInfo = ""
							setupFileForm(ctx, form, file, r, service, downloadFolder)
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
