package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"string_um/string/main/funcs"
	"string_um/string/main/tui/globals"
)

var loginForm = tview.NewForm().SetButtonsAlign(tview.AlignCenter)

func login() {
	InfoBoxInstance.SetText("")
	password := loginForm.GetFormItemByLabel("Password: ").(*tview.InputField).GetText()
	if err := funcs.Login(password); err != nil {
		errorMssg := err.Error()
		globals.LowerTextView.SetText(errorMssg).SetTextColor(tcell.ColorRed)
		return
	} else {
		globals.LoginSuccessChan <- true
		globals.Pages.SwitchToPage("main")
	}
}

func LoginForm() *tview.Form {
	loginForm.SetTitleAlign(tview.AlignCenter)
	loginForm.AddPasswordField("Password: ", "", 30, '*', func(text string) {
		globals.LowerTextView.SetText("")
	})
	loginForm.AddButton("Login", login)
	loginForm.SetLabelColor(tcell.NewRGBColor(228, 179, 99)).SetFieldBackgroundColor(tcell.NewRGBColor(224, 223, 213)).SetFieldTextColor(tcell.NewRGBColor(50, 55, 57)).SetButtonBackgroundColor(tcell.NewRGBColor(228, 179, 99)).SetButtonTextColor(tcell.ColorBlack)

	return loginForm
}
