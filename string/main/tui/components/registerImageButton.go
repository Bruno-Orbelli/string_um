package components

import (
	"string_um/string/main/funcs"
	"string_um/string/main/tui/globals"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func registerSecondFactor() {
	capturedEncoding, err := funcs.CaptureMultipleImagesEncoding(15)
	if err != nil {
		errorMssg := err.Error()
		globals.LowerTextView.SetText(errorMssg).SetTextColor(tcell.ColorRed)
		return
	}
	if err := funcs.Register(password, capturedEncoding); err != nil {
		errorMssg := err.Error()
		globals.LowerTextView.SetText(errorMssg).SetTextColor(tcell.ColorRed)
		return
	} else {
		globals.LowerTextView.SetText("Registration successful!").SetTextColor(tcell.ColorGreen)
		password = ""
		globals.Pages.SwitchToPage("login")
	}
}

func RegisterImageButton() *tview.Button {
	registerImageButton := tview.NewButton("Take pictures")
	registerImageButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.NewRGBColor(228, 179, 99)))
	registerImageButton.SetSelectedFunc(registerSecondFactor)
	registerImageButton.SetLabelColor(tcell.ColorBlack)
	registerImageButton.SetBackgroundColor(tcell.NewRGBColor(228, 179, 99))
	registerImageButton.SetLabelColorActivated(tcell.NewRGBColor(228, 179, 99))
	registerImageButton.SetBackgroundColorActivated(tcell.ColorBlack)
	registerImageButton.SetTitleAlign(tview.AlignCenter)

	return registerImageButton
}
