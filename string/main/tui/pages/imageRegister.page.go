package pages

import (
	"string_um/string/main/tui/components"
	"string_um/string/main/tui/globals"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func BuildImageRegisterPage() tview.Primitive {
	flex := tview.NewFlex()
	flex.SetBorder(true).SetBorderColor(tcell.NewRGBColor(228, 179, 99))

	flex1 := tview.NewFlex()
	flex2 := tview.NewFlex()
	flex3 := tview.NewFlex()
	flex4 := tview.NewFlex()
	flex5 := tview.NewFlex()
	flex6 := tview.NewFlex()

	logo := components.Logo(tcell.NewRGBColor(228, 179, 99))

	subtitle := components.RegisterImageText()

	button := components.RegisterImageButton()

	flex2.SetDirection(tview.FlexRow).AddItem(logo, 0, 1, true)

	flex3.SetDirection(tview.FlexRow).AddItem(subtitle, 0, 1, true)

	flex5.SetDirection(tview.FlexRow).AddItem(tview.NewBox(), 0, 3, true)
	flex5.SetDirection(tview.FlexRow).AddItem(button, 0, 1, true)
	flex5.SetDirection(tview.FlexRow).AddItem(tview.NewBox(), 0, 3, true)

	flex4.SetDirection(tview.FlexColumn).AddItem(tview.NewBox(), 0, 4, true)
	flex4.SetDirection(tview.FlexColumn).AddItem(flex5, 0, 1, true)
	flex4.SetDirection(tview.FlexColumn).AddItem(tview.NewBox(), 0, 4, true)

	flex6.SetDirection(tview.FlexRow).AddItem(globals.LowerTextView, 0, 1, true)

	flex.SetDirection(tview.FlexRow).AddItem(flex1, 0, 1, true)
	flex.SetDirection(tview.FlexRow).AddItem(flex2, 0, 2, true)
	flex.SetDirection(tview.FlexRow).AddItem(flex3, 0, 1, true)
	flex.SetDirection(tview.FlexRow).AddItem(flex4, 0, 1, true)
	flex.SetDirection(tview.FlexRow).AddItem(flex6, 0, 1, true)

	return flex
}
