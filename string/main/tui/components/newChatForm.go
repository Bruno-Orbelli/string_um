package components

import (
	"string_um/string/main/funcs"
	"string_um/string/main/tui/globals"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var newChatForm = tview.NewForm()

func openChat() {
	contactName := newChatForm.GetFormItemByLabel("Contact name: ").(*tview.InputField).GetText()
	if contactName == "" {
		InfoBoxInstance.Clear()
		UpdateInfo(true, "Contact name can't be empty.")
		return
	}
	contact, err := funcs.GetContactByName(contactName)
	if err != nil {
		InfoBoxInstance.Clear()
		UpdateInfo(true, err.Error())
		return
	} else if contact == nil {
		InfoBoxInstance.Clear()
		UpdateInfo(true, "Contact not found.")
		return
	}
	_, err = funcs.GetChatWithContact(contact.ID)
	if err != nil {
		InfoBoxInstance.Clear()
		UpdateInfo(true, err.Error())
		return
	}
	globals.ChatsRefreshedChan <- true
	InfoBoxInstance.Clear()
	UpdateInfo(false, "Chat opened.")
}

func NewChatForm() *tview.Form {
	globals.LowerTextView.SetText("")
	InfoBoxInstance.Clear()
	newChatForm.SetBorder(true)
	newChatForm.SetBorderColor(tview.Styles.ContrastBackgroundColor)
	newChatForm.SetTitle("Open a new chat")
	newChatForm.SetTitleAlign(tview.AlignCenter)
	newChatForm.SetTitleColor(tcell.NewRGBColor(224, 223, 213))
	newChatForm.SetLabelColor(tcell.NewRGBColor(224, 223, 213))
	newChatForm.SetFieldBackgroundColor(tcell.NewRGBColor(224, 223, 213))
	newChatForm.SetFieldTextColor(tcell.NewRGBColor(50, 55, 57))
	newChatForm.SetButtonBackgroundColor(tview.Styles.ContrastBackgroundColor)
	newChatForm.SetButtonsAlign(tview.AlignCenter)
	newChatForm.SetCancelFunc(goBack)
	newChatForm.AddInputField("Contact name: ", "", 40, nil, nil)
	newChatForm.AddButton("Open", openChat)
	return newChatForm
}
