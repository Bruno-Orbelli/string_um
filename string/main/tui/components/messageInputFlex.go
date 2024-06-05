package components

import (
	"string_um/string/main/funcs"
	"string_um/string/main/tui/globals"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

var messageTextBox = tview.NewForm()
var currentChatID uuid.UUID

func sendMessage() {
	message := messageTextBox.GetFormItemByLabel("").(*tview.InputField).GetText()
	ownUser, _, err := funcs.GetOwnUser()
	if err != nil {
		panic(err)
	}
	if message != "" {
		funcs.AddMessageToBeSent(currentChatID, ownUser.ID, message)
		globals.MessagesRefreshedChan <- true
		messageTextBox.GetFormItemByLabel("").(*tview.InputField).SetText("")
	}
}

func MessageInputFlex(chatID uuid.UUID) *tview.Flex {
	currentChatID = chatID
	messageInputFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	auxFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	auxFlex2 := tview.NewFlex().SetDirection(tview.FlexColumn)
	auxFlex2.AddItem(tview.NewBox(), 0, 1, false)
	auxFlex2.AddItem(SendMessageButton(), 0, 4, false)
	auxFlex2.AddItem(tview.NewBox(), 0, 1, false)
	auxFlex.AddItem(tview.NewBox(), 0, 1, false)
	auxFlex.AddItem(auxFlex2, 0, 1, false)
	auxFlex.AddItem(tview.NewBox(), 0, 1, false)
	messageInputFlex.AddItem(MessageTextBox(), 0, 5, false)
	messageInputFlex.AddItem(auxFlex, 0, 1, false)

	return messageInputFlex
}

func SendMessageButton() *tview.Button {
	sendMessageButton := tview.NewButton("Send")
	sendMessageButton.SetStyle(tcell.StyleDefault.Foreground(tcell.NewRGBColor(232, 233, 235)).Background(tcell.NewRGBColor(2, 128, 0)))
	sendMessageButton.SetSelectedFunc(sendMessage)
	sendMessageButton.SetLabelColor(tcell.NewRGBColor(232, 233, 235))
	sendMessageButton.SetBackgroundColor(tcell.NewRGBColor(2, 128, 0))
	sendMessageButton.SetLabelColorActivated(tcell.NewRGBColor(2, 128, 0))
	sendMessageButton.SetBackgroundColorActivated(tcell.NewRGBColor(232, 233, 235))

	return sendMessageButton
}

func MessageTextBox() *tview.Form {
	messageTextBox.SetButtonsAlign(tview.AlignCenter)
	messageTextBox.AddInputField("", "", 70, nil, nil)
	messageTextBox.SetLabelColor(tcell.NewRGBColor(228, 179, 99))
	messageTextBox.SetFieldBackgroundColor(tcell.NewRGBColor(224, 223, 213))
	messageTextBox.SetFieldTextColor(tcell.NewRGBColor(50, 55, 57))

	return messageTextBox
}
