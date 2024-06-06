package components

import (
	"string_um/string/main/tui/globals"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var registerForm = tview.NewForm().SetButtonsAlign(tview.AlignCenter)
var password string

func checkPassword() bool {
	inputedPassword := registerForm.GetFormItemByLabel("Password: ").(*tview.InputField).GetText()
	confirmedPassword := registerForm.GetFormItemByLabel("Confirm password: ").(*tview.InputField).GetText()
	return inputedPassword == confirmedPassword
}

func registerFirstFactor() {
	inputedPassword := registerForm.GetFormItemByLabel("Password: ").(*tview.InputField).GetText()
	if !checkPassword() {
		globals.LowerTextView.SetText("Passwords do not match.").SetTextColor(tcell.ColorRed)
		return
	}
	password = inputedPassword
	globals.Pages.SwitchToPage("imageRegister")
}

func RegisterForm() *tview.Form {
	registerForm.SetTitleAlign(tview.AlignCenter)
	registerForm.AddPasswordField("Password: ", "", 30, '*', func(text string) {
		globals.LowerTextView.SetText("")
	})
	registerForm.AddPasswordField("Confirm password: ", "", 30, '*', func(text string) {
		globals.LowerTextView.SetText("")
	})
	registerForm.AddButton("Register", registerFirstFactor)
	registerForm.SetLabelColor(tcell.NewRGBColor(228, 179, 99)).SetFieldBackgroundColor(tcell.NewRGBColor(224, 223, 213)).SetFieldTextColor(tcell.NewRGBColor(50, 55, 57)).SetButtonBackgroundColor(tcell.NewRGBColor(228, 179, 99)).SetButtonTextColor(tcell.ColorBlack)

	return registerForm
}
