package tui

import (
	"context"
	"time"

	"github.com/rivo/tview"
	"github.com/vukit/gophkeeper/internal/client"
	"github.com/vukit/gophkeeper/internal/client/model"
)

// Cards компонент реализует текстовый интерфейс CRUD для банковских карт
func (r *TUI) Cards(ctx context.Context, user *model.User, service client.GophKeeperService) *tview.Flex {
	card := &model.Card{}

	layout := tview.NewFlex()

	list := tview.NewList()
	list.SetTitle("[ Cards ]").SetBorder(true).SetBorderPadding(1, 1, 1, 1)
	list.ShowSecondaryText(false)

	form := tview.NewForm()
	card.ID = 0
	card.Bank = ""
	card.Number = ""
	card.Date = ""
	card.CVV = ""
	card.MetaInfo = ""
	setupCardForm(ctx, form, card, r, service)

	layout.AddItem(list, 35, 0, false).AddItem(form, 0, 1, true)

	go cardsUpdateList(ctx, card, list, form, r, service)

	return layout
}

func setupCardForm(
	ctx context.Context,
	form *tview.Form,
	card *model.Card,
	r *TUI,
	service client.GophKeeperService,
) {
	form.Clear(true)
	form.
		AddInputField("Bank", card.Bank, 30, nil, func(text string) { card.Bank = text }).
		AddInputField("Number", card.Number, 20, nil, func(text string) { card.Number = text }).
		AddInputField("Date", card.Date, 6, nil, func(text string) { card.Date = text }).
		AddInputField("CVV", card.CVV, 4, nil, func(text string) { card.CVV = text }).
		AddInputField("Metainfo", card.MetaInfo, 30, nil, func(text string) { card.MetaInfo = text }).
		AddButton("Save", func() {
			err := card.Validate()
			if err != nil {
				r.alertChannel <- err.Error()

				return
			}

			err = service.SaveCard(ctx, card)
			if err != nil {
				r.alertChannel <- err.Error()

				return
			}

			r.alertChannel <- "card was successfully saved"
			card.ID = 0
			card.Bank = ""
			card.Number = ""
			card.Date = ""
			card.CVV = ""
			card.MetaInfo = ""
			setupCardForm(ctx, form, card, r, service)
		}).
		AddButton("Cancel", func() {
			card.ID = 0
			card.Bank = ""
			card.Number = ""
			card.Date = ""
			card.CVV = ""
			card.MetaInfo = ""
			setupCardForm(ctx, form, card, r, service)
		})

	if idx := form.GetButtonIndex("Delete"); idx != -1 {
		form.RemoveButton(idx)
	}

	form.SetBorder(true).SetTitle("[ New card ]").SetTitleAlign(tview.AlignLeft).SetBorderPadding(1, 0, 2, 0)

	r.app.SetFocus(form)
}

func cardsUpdateList(
	ctx context.Context,
	card *model.Card,
	list *tview.List,
	form *tview.Form,
	r *TUI,
	service client.GophKeeperService,
) {
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			cards, err := service.GetCards(ctx)
			if err != nil {
				r.alertChannel <- err.Error()

				continue
			}

			currentItemIndex := list.GetCurrentItem()

			list.Clear()

			for _, currentcard := range cards {
				currentCard := currentcard
				list.AddItem(currentCard.MetaInfo, "", rune(0), func() {
					card = &currentCard
					setupCardForm(ctx, form, card, r, service)
					form.SetTitle("[ Edit card ]")
					if idx := form.GetButtonIndex("Delete"); idx == -1 {
						form.AddButton("Delete", func() {
							errService := service.DeleteCard(ctx, card)
							if errService != nil {
								r.alertChannel <- errService.Error()

								return
							}
							r.alertChannel <- "card was successfully deleted"
							card.ID = 0
							card.Bank = ""
							card.Number = ""
							card.Date = ""
							card.CVV = ""
							card.MetaInfo = ""
							setupCardForm(ctx, form, card, r, service)
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
