package pages

import (
	"os"
	"string_um/string/main/funcs"
	"string_um/string/models"

	"string_um/string/main/tui/components"
	"string_um/string/main/tui/globals"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

var chats []models.ChatDTO
var chatListView = components.ChatList()
var selectedChatList = components.SelectedChatList()
var selectedChatID *uuid.UUID

func updateChatList() {
	chats = getChats()
	chatListView.Clear()
	for _, chat := range chats {
		chatListView.AddItem(chat.ContactName, chat.ContactID, 0, nil)
	}
	chatListView.SetCurrentItem(0)
}

func mainLoop(sigChan chan os.Signal, app *tview.Application) {
	for {
		select {
		case <-globals.ChatsReadyChan:
			app.QueueUpdateDraw(updateChatList)
			funcs.AddContactAddressesForUnknownContacts(components.Libp2pHost)
		case <-globals.ChatsRefreshedChan:
			app.QueueUpdateDraw(updateChatList)
			funcs.AddContactAddressesForUnknownContacts(components.Libp2pHost)
		case <-globals.MessagesRefreshedChan:
			chats = getChats()
			if selectedChatID != nil {
				app.QueueUpdateDraw(displayMessages)
			}
		case <-sigChan:
			return
		}
	}
}

func getChats() []models.ChatDTO {
	chats, err := funcs.GetChatsAndInfo()
	if err != nil {
		panic(err)
	}
	if chats == nil {
		return []models.ChatDTO{}
	}
	return chats
}

func getSelectedChat(chatID uuid.UUID) *models.ChatDTO {
	for _, chat := range chats {
		if chat.ID == chatID {
			return &chat
		}
	}
	return nil
}

func displayMessages() {
	selectedChat := getSelectedChat(*selectedChatID)
	messages := selectedChat.Messages
	ownUser, _, err := funcs.GetOwnUser()
	if err != nil {
		panic(err)
	}
	selectedChatList.Clear()
	for _, message := range messages {
		if message.SentByID == ownUser.ID {
			selectedChatList.AddItem(
				components.Message(true, "You", message.Message),
				message.SentAt.Format("02 Jan 2006 15:04"),
				0,
				nil,
			)
		} else {
			selectedChatList.AddItem(
				components.Message(false, selectedChat.ContactName, message.Message),
				message.SentAt.Format("02 Jan 2006 15:04"),
				0,
				nil,
			)
		}
	}
	selectedChatList.SetCurrentItem(0)
}

func BuildMainPage(app *tview.Application, sigChan chan os.Signal) tview.Primitive {
	flex := tview.NewFlex()
	flex.SetBorder(false).SetBorderAttributes(tcell.AttrDim).SetBorderColor(tcell.NewRGBColor(228, 179, 99))

	flex1 := tview.NewFlex()
	flex2 := tview.NewFlex()

	upperInfo := components.InfoText()
	logo := components.Logo(tcell.NewRGBColor(232, 233, 235))

	for _, chat := range chats {
		chatListView.AddItem(chat.ContactName, chat.ContactID, 0, nil)
	}

	// Central component
	flex1.SetBorder(true).SetBorderColor(tcell.NewRGBColor(232, 233, 235))
	flex1.SetDirection(tview.FlexRow).AddItem(upperInfo, 0, 1, true)
	flex1.SetDirection(tview.FlexRow).AddItem(logo, 0, 4, true)
	flex1.SetDirection(tview.FlexRow).AddItem(tview.NewBox(), 0, 1, true)

	// Select chat
	var chatTitleView *tview.TextView
	chatListView.SetSelectedFunc(
		func(row int, mainText, secondaryText string, shortcut rune) {
			selectedChat := getSelectedChat(chats[row].ID)
			selectedChatID = &selectedChat.ID
			chatTitleView = components.ChatTitle(selectedChat.ContactName)
			selectedChatList.Clear()
			flex1.Clear()
			flex1.AddItem(chatTitleView, 0, 1, true)
			displayMessages()
			flex1.AddItem(selectedChatList, 0, 7, true)
			flex1.AddItem(components.MessageInputFlex(selectedChat.ID), 0, 1, false)
		},
	)

	// Option menu
	optionMenu := components.OptionsMenu(sigChan, flex1)

	// Right-side bar
	flex2.SetDirection(tview.FlexRow).AddItem(optionMenu, 0, 5, true)
	flex2.SetDirection(tview.FlexRow).AddItem(components.InfoBoxInstance, 0, 2, true)

	flex.SetDirection(tview.FlexColumn).AddItem(chatListView, 0, 1, true)
	flex.SetDirection(tview.FlexColumn).AddItem(flex1, 0, 4, true)
	flex.SetDirection(tview.FlexColumn).AddItem(flex2, 0, 1, true)

	go mainLoop(sigChan, app)

	return flex
}
