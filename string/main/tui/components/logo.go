package components

import (
	"string_um/string/main/tui/globals"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Logo(color tcell.Color) *tview.TextView {
	logo := tview.NewTextView()
	logo.Write(globals.LoadTitle())
	logo.SetTextAlign(tview.AlignCenter)
	logo.SetTextColor(color)

	return logo
}
