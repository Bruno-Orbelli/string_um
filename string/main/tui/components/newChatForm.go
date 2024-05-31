package components

import (
	"string_um/string/main/funcs"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var newChatForm = tview.NewForm()

func openChat() {
	contactName := newChatForm.GetFormItemByLabel("Contact name: ").(*tview.InputField).GetText()
	if contactName == "" {
		LowerTextView.SetText("Contact name can't be empty.").SetTextColor(tcell.ColorRed)
		return
	}
	contact, err := funcs.GetContactByName(contactName)
	if err != nil {
		LowerTextView.SetText(err.Error()).SetTextColor(tcell.ColorRed)
		return
	} else if contact == nil {
		LowerTextView.SetText("Contact not found.").SetTextColor(tcell.ColorRed)
		return
	}
	_, err = funcs.GetChatWithContact(contact.ID)
	if err != nil {
		LowerTextView.SetText(err.Error())
		return
	}
	ChatsRefreshedChan <- true
	LowerTextView.SetText("Chat opened.").SetTextColor(tcell.ColorGreen)
}

func NewChatForm() *tview.Form {
	LowerTextView.SetText("")
	newChatForm.SetBorder(true)
	newChatForm.SetTitle("Open a new chat")
	newChatForm.SetTitleAlign(tview.AlignLeft)
	newChatForm.SetLabelColor(tview.Styles.PrimaryTextColor)
	newChatForm.SetFieldBackgroundColor(tview.Styles.ContrastBackgroundColor)
	newChatForm.SetFieldTextColor(tview.Styles.PrimaryTextColor)
	newChatForm.SetButtonBackgroundColor(tview.Styles.ContrastBackgroundColor)
	newChatForm.SetButtonTextColor(tview.Styles.PrimaryTextColor)
	newChatForm.SetButtonsAlign(tview.AlignCenter)
	newChatForm.SetButtonsAlign(tview.AlignCenter)
	newChatForm.SetLabelColor(tview.Styles.PrimaryTextColor)
	newChatForm.SetCancelFunc(goBack)
	newChatForm.AddInputField("Contact name: ", "", 40, nil, nil)
	newChatForm.AddButton("Open", openChat)
	newChatForm.AddFormItem(LowerTextView)
	return newChatForm
}
