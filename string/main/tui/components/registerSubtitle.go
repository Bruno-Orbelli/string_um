package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RegisterSubtitle() *tview.TextView {
	registerSubtitle := tview.NewTextView()
	registerSubtitle.SetText("Welcome to String! To get started, please create a strong password for your account.\nThis password will be used to encrypt your private key and messages.")
	registerSubtitle.SetTextAlign(tview.AlignCenter)
	registerSubtitle.SetTextColor(tcell.NewRGBColor(232, 233, 235))

	return registerSubtitle
}
