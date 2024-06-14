package components

import (
	"string_um/string/main/tui/globals"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var newOwnInfoForm = tview.NewForm()

func getOwnInfo() {
	<-globals.OwnInfoAvailableChan
	newOwnInfoForm.AddInputField("Own multihash: ", globals.OwnUserHash, 46, nil, nil)
}

func copyOwnMultihash() {
	err := clipboard.WriteAll(globals.OwnUserHash)
	if err != nil {
		UpdateInfo(true, "Error copying to clipboard: "+err.Error())
		return
	}
	UpdateInfo(false, "Copied to clipboard.")
}

func OwnInfoForm() *tview.Form {
	globals.LowerTextView.SetText("")
	InfoBoxInstance.Clear()
	newOwnInfoForm.SetBorder(true)
	newOwnInfoForm.SetTitle("My info")
	newOwnInfoForm.SetTitleAlign(tview.AlignCenter)
	newOwnInfoForm.SetTitleColor(tcell.NewRGBColor(232, 233, 235))
	newOwnInfoForm.SetBorderColor(tcell.NewRGBColor(202, 21, 202))
	newOwnInfoForm.SetBackgroundColor(tcell.ColorBlack)
	newOwnInfoForm.SetLabelColor(tcell.NewRGBColor(232, 233, 235))
	newOwnInfoForm.SetFieldBackgroundColor(tcell.NewRGBColor(224, 223, 213))
	newOwnInfoForm.SetFieldTextColor(tcell.NewRGBColor(50, 55, 57))
	newOwnInfoForm.SetButtonBackgroundColor(tcell.NewRGBColor(202, 21, 202))
	newOwnInfoForm.SetButtonTextColor(tcell.NewRGBColor(224, 223, 213))
	newOwnInfoForm.SetButtonsAlign(tview.AlignCenter)
	newOwnInfoForm.SetCancelFunc(goBack)
	newOwnInfoForm.AddButton("Copy to clipboard", copyOwnMultihash)

	/* ownMultihash := tview.NewInputField()
	ownMultihash.SetLabel("Own multihash: ")
	ownMultihash.SetFieldWidth(40)
	ownMultihash.SetFieldBackgroundColor(tcell.NewRGBColor(224, 223, 213))
	ownMultihash.SetFieldTextColor(tcell.NewRGBColor(50, 55, 57))
	ownMultihash.SetLabelColor(tcell.NewRGBColor(232, 233, 235))
	ownMultihash.SetLabelWidth(20)
	ownMultihash.SetDisabled(true) */

	go getOwnInfo()

	return newOwnInfoForm
}
