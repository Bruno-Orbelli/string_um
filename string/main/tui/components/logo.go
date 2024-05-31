package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Logo(color tcell.Color) *tview.TextView {
	logo := tview.NewTextView()
	logo.Write(LoadTitle())
	logo.SetTextAlign(tview.AlignCenter)
	logo.SetTextColor(color)

	return logo
}
