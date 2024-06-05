package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SelectedChatList() *tview.List {
	selectedChat := tview.NewList()
	selectedChat.ShowSecondaryText(true)
	selectedChat.SetSecondaryTextColor(tcell.NewRGBColor(111, 115, 116))
	selectedChat.SetMainTextColor(tcell.NewRGBColor(232, 233, 235))
	selectedChat.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	selectedChat.SetSelectedTextColor(tcell.NewRGBColor(232, 233, 235))
	selectedChat.SetSelectedBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	selectedChat.SetBorder(false)

	return selectedChat
}
