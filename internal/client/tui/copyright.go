package tui

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

// nolint
var (
	buildVersion string = "dev"
	buildData    string = time.Now().Format("2006-01-02")
)

// Copyright выводит название, версию и дату сборки программы
func Copyright(align int) *tview.TextView {
	text := fmt.Sprintf("GophKeeper %s/%s", buildVersion, buildData)

	return tview.NewTextView().SetText(text).SetTextAlign(align)
}
