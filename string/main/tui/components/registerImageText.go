package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RegisterImageText() *tview.TextView {
	registerImageText := tview.NewTextView()
	registerImageText.SetDynamicColors(true)
	registerImageText.SetRegions(true)
	registerImageText.SetBorder(false)
	registerImageText.SetTextAlign(tview.AlignCenter)
	registerImageText.SetTextColor(tcell.NewRGBColor(232, 233, 235))
	registerImageText.SetText("To finish the registration process, we will capture a quick succession of images of your face.\nDon't worry, we will not store these images, they will be used to create a unique encoding for authentication purposes.\nMake sure you are in a well-lit area, your webcam cover is removed, and you are looking directly at the camera.\n\nPress the button below when you are ready.")
	return registerImageText
}
