package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InfoText() *tview.TextView {
	upperInfo := tview.NewTextView()
	upperInfo.SetText("String - Secure Messaging.\tVersion 1.0.0")
	upperInfo.SetTextStyle(tcell.StyleDefault.Italic(true))
	upperInfo.SetTextAlign(tview.AlignCenter)
	upperInfo.SetTextColor(tcell.NewRGBColor(232, 233, 235))

	return upperInfo
}
