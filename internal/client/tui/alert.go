package tui

import (
	"context"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Alert выводит сообщение в заголовке интерфейса
func Alert(ctx context.Context, top, bottom, left, right, align int, clearInterval int64) (alert *tview.TextView, ch chan string) {
	ch = make(chan string, 5)

	alert = tview.NewTextView().SetTextColor(tcell.ColorYellow)
	alert.SetBorderPadding(top, bottom, left, right)
	alert.SetTextAlign(align)

	if clearInterval != 0 {
		go func() {
			interval := time.Duration(clearInterval) * time.Second
			timer := time.NewTimer(0)

			for {
				select {
				case <-ctx.Done():
					timer.Stop()

					return
				case msg := <-ch:
					alert.SetText(msg)
					timer.Reset(interval)
				case <-timer.C:
					alert.SetText("")
				}
			}
		}()
	}

	return alert, ch
}
