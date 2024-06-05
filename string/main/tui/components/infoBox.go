package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var InfoBoxInstance = InfoBox()

func UpdateInfo(isError bool, message string) {
	InfoBoxInstance.SetText("")
	if isError {
		InfoBoxInstance.SetTitleColor(tcell.ColorRed)
		InfoBoxInstance.SetTextColor(tcell.ColorRed)
		InfoBoxInstance.SetBorderColor(tcell.ColorRed)
	} else {
		InfoBoxInstance.SetTitleColor(tcell.ColorGreen)
		InfoBoxInstance.SetTextColor(tcell.ColorGreen)
		InfoBoxInstance.SetBorderColor(tcell.ColorGreen)
	}
	InfoBoxInstance.SetText(message)
}

func InfoBox() *tview.TextView {
	infoBoxInstance := tview.NewTextView()
	infoBoxInstance.SetTextAlign(tview.AlignCenter)
	infoBoxInstance.SetBorder(true)
	infoBoxInstance.SetBorderColor(tcell.NewRGBColor(50, 55, 57))
	infoBoxInstance.SetTitle("Info")
	infoBoxInstance.SetTitleColor(tcell.NewRGBColor(50, 55, 57))
	infoBoxInstance.SetTitleAlign(tview.AlignCenter)

	return infoBoxInstance
}
