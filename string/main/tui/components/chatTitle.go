package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ChatTitle(contactName string) *tview.TextView {
	chatTitle := tview.NewTextView()
	chatTitle.SetText(contactName)
	chatTitle.SetTextColor(tcell.NewRGBColor(232, 233, 235))
	chatTitle.SetTextAlign(tview.AlignCenter)

	return chatTitle
}
