package pages

import (
	funcs "string_um/string/main/funcs"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var confirmedPassword = ""

func checkPassword() bool {
	return inputedPassword == confirmedPassword
}

func register() {
	if !checkPassword() {
		lowerTextView.SetText("Passwords do not match.").SetTextColor(tcell.ColorRed)
		return
	}
	if err := funcs.Register(inputedPassword); err != nil {
		errorMssg := err.Error()
		lowerTextView.SetText(errorMssg).SetTextColor(tcell.ColorRed)
		return
	} else {
		inputedPassword = ""
		confirmedPassword = ""
		lowerTextView.SetText("Registration successful!").SetTextColor(tcell.ColorGreen)
		Pages.SwitchToPage("login")
	}
}

func BuildRegisterPage() tview.Primitive {
	flex := tview.NewFlex()
	flex.SetBorder(true).SetBorderColor(tcell.NewRGBColor(228, 179, 99))

	flex1 := tview.NewFlex()
	flex2 := tview.NewFlex()
	flex3 := tview.NewFlex()
	flex4 := tview.NewFlex()
	flex5 := tview.NewFlex()

	logo := tview.NewTextView()
	logo.Write(loadTitle())
	logo.SetTextAlign(tview.AlignCenter)
	logo.SetTextColor(tcell.NewRGBColor(228, 179, 99))

	subtitle := tview.NewTextView()
	subtitle.SetText("Welcome to String! To get started, please create a strong password for your account.\nThis password will be used to encrypt your private key and messages.")
	subtitle.SetTextAlign(tview.AlignCenter)
	subtitle.SetTextColor(tcell.NewRGBColor(232, 233, 235))

	form := tview.NewForm().SetButtonsAlign(tview.AlignCenter)
	form.SetTitleAlign(tview.AlignCenter)
	form.AddPasswordField("Password: ", "", 30, '*', func(text string) {
		lowerTextView.SetText("")
		inputedPassword = text
	})
	form.AddPasswordField("Confirm password: ", "", 30, '*', func(text string) {
		lowerTextView.SetText("")
		confirmedPassword = text
	})
	form.AddButton("Register", register)
	form.SetLabelColor(tcell.NewRGBColor(228, 179, 99)).SetFieldBackgroundColor(tcell.NewRGBColor(224, 223, 213)).SetFieldTextColor(tcell.NewRGBColor(50, 55, 57)).SetButtonBackgroundColor(tcell.NewRGBColor(228, 179, 99)).SetButtonTextColor(tcell.ColorBlack)

	flex2.SetDirection(tview.FlexRow).AddItem(logo, 0, 1, true)

	flex3.SetDirection(tview.FlexRow).AddItem(subtitle, 0, 1, true)

	flex4.SetDirection(tview.FlexColumn).AddItem(tview.NewBox(), 0, 1, true)
	flex4.SetDirection(tview.FlexColumn).AddItem(form, 0, 1, true)
	flex4.SetDirection(tview.FlexColumn).AddItem(tview.NewBox(), 0, 1, true)

	flex5.SetDirection(tview.FlexRow).AddItem(lowerTextView, 0, 1, true)

	flex.SetDirection(tview.FlexRow).AddItem(flex1, 0, 1, true)
	flex.SetDirection(tview.FlexRow).AddItem(flex2, 0, 2, true)
	flex.SetDirection(tview.FlexRow).AddItem(flex3, 0, 1, true)
	flex.SetDirection(tview.FlexRow).AddItem(flex4, 0, 2, true)
	flex.SetDirection(tview.FlexRow).AddItem(flex5, 0, 1, true)

	return flex
}
