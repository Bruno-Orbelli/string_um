package pages

import (
	funcs "string_um/string/main/funcs"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func login() {
	lowerTextView.SetText("")
	if err := funcs.Login(inputedPassword); err != nil {
		errorMssg := err.Error()
		lowerTextView.SetText(errorMssg).SetTextColor(tcell.ColorRed)
		return
	} else {
		inputedPassword = ""
		Pages.SwitchToPage("main")
	}
}

func BuildLoginPage() tview.Primitive {
	flex := tview.NewFlex()
	flex.SetBorder(true).SetBorderColor(tcell.NewRGBColor(228, 179, 99))

	flex1 := tview.NewFlex()
	flex2 := tview.NewFlex()
	flex3 := tview.NewFlex()
	flex4 := tview.NewFlex()

	logo := tview.NewTextView()
	logo.Write(loadTitle())
	logo.SetTextAlign(tview.AlignCenter)
	logo.SetTextColor(tcell.NewRGBColor(228, 179, 99))

	form := tview.NewForm().SetButtonsAlign(tview.AlignCenter)
	form.SetTitleAlign(tview.AlignCenter)
	form.AddPasswordField("Password: ", "", 30, '*', func(text string) {
		lowerTextView.SetText("")
		inputedPassword = text
	})
	form.AddButton("Login", login)
	form.SetLabelColor(tcell.NewRGBColor(228, 179, 99)).SetFieldBackgroundColor(tcell.NewRGBColor(224, 223, 213)).SetFieldTextColor(tcell.NewRGBColor(50, 55, 57)).SetButtonBackgroundColor(tcell.NewRGBColor(228, 179, 99)).SetButtonTextColor(tcell.ColorBlack)

	flex1.SetDirection(tview.FlexRow)

	flex2.SetDirection(tview.FlexRow).AddItem(logo, 0, 1, true)

	flex3.SetDirection(tview.FlexColumn).AddItem(tview.NewBox(), 0, 1, true)
	flex3.SetDirection(tview.FlexColumn).AddItem(form, 0, 1, true)
	flex3.SetDirection(tview.FlexColumn).AddItem(tview.NewBox(), 0, 1, true)

	flex4.SetDirection(tview.FlexColumn).AddItem(lowerTextView, 0, 1, true)

	flex.SetDirection(tview.FlexRow).AddItem(flex1, 0, 1, true)
	flex.SetDirection(tview.FlexRow).AddItem(flex2, 0, 1, true)
	flex.SetDirection(tview.FlexRow).AddItem(flex3, 0, 1, true)
	flex.SetDirection(tview.FlexRow).AddItem(flex4, 0, 1, true)

	return flex
}
