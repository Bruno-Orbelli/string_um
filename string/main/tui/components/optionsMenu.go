package components

import (
	"os"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var sigChannel chan os.Signal
var containingFlex *tview.Flex
var contactForm = AddContactForm()
var chatForm = NewChatForm()
var ownInfoForm = OwnInfoForm()

func onNewChatSelected() {
	containingFlex.Clear()
	containingFlex.AddItem(InfoText(), 0, 1, false)
	containingFlex.AddItem(chatForm, 0, 5, true)
}

func onAddNewContactSelected() {
	containingFlex.Clear()
	containingFlex.AddItem(InfoText(), 0, 1, false)
	containingFlex.AddItem(contactForm, 0, 5, true)
}

func onCheckMyInfoSelected() {
	containingFlex.Clear()
	containingFlex.AddItem(InfoText(), 0, 1, false)
	containingFlex.AddItem(ownInfoForm, 0, 5, true)
}

func onQuitSelected() {
	sigChannel <- syscall.SIGINT
}

func OptionsMenu(sigChan chan os.Signal, flex *tview.Flex) *tview.List {
	sigChannel = sigChan
	containingFlex = flex

	// Option menu
	optionMenu := tview.NewList()
	optionMenu.ShowSecondaryText(false)
	optionMenu.SetMainTextColor(tcell.NewRGBColor(0, 209, 202))
	optionMenu.SetSecondaryTextColor(tcell.NewRGBColor(0, 209, 202))
	optionMenu.SetSelectedTextColor(tcell.ColorBlack)
	optionMenu.SetSelectedBackgroundColor(tcell.NewRGBColor(0, 209, 202))
	optionMenu.SetBorder(true).SetBorderColor(tcell.NewRGBColor(0, 209, 202))
	optionMenu.SetTitle("Menu").SetTitleColor(tcell.NewRGBColor(0, 209, 202))

	optionMenu.AddItem("New chat", "Create a new chat", 0, onNewChatSelected)
	optionMenu.AddItem("Add new contact", "Add a new contact", 0, onAddNewContactSelected)
	// optionMenu.AddItem("Check my contacts", "Check your contacts", 0, nil)
	optionMenu.AddItem("Check my info", "Check your personal info", 0, onCheckMyInfoSelected)
	optionMenu.AddItem("Quit", "Quit the application", 0, onQuitSelected)

	return optionMenu
}
