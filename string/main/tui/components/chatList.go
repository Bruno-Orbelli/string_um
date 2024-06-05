package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ChatList() *tview.List {
	chatList := tview.NewList()
	chatList.ShowSecondaryText(true)
	chatList.SetMainTextColor(tcell.NewRGBColor(232, 233, 235))
	chatList.SetSecondaryTextColor(tcell.NewRGBColor(241, 217, 177))
	chatList.SetSelectedTextColor(tcell.ColorBlack)
	chatList.SetSelectedBackgroundColor(tcell.NewRGBColor(228, 179, 99))
	chatList.SetBorder(true).SetBorderColor(tcell.NewRGBColor(182, 143, 79))
	chatList.SetTitle("Chats").SetTitleColor(tcell.NewRGBColor(228, 179, 99))

	return chatList
}
