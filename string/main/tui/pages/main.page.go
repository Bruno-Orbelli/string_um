package pages

import (
	"fmt"
	"os"
	"string_um/string/main/funcs"
	"string_um/string/models"

	"string_um/string/main/tui/components"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

var chats []models.ChatDTO
var chatListView = components.ChatList()

// TODO: Add a chan to refresh the chats

func updateChatList() {
	chats = getChats()
	chatListView.Clear()
	for _, chat := range chats {
		chatListView.AddItem(chat.ContactName, chat.ContactID, 0, nil)
	}
	chatListView.SetCurrentItem(0)
}

func waitForChats(sigChan chan os.Signal, app *tview.Application) {
	for {
		select {
		case <-components.ChatsReadyChan:
			app.QueueUpdateDraw(updateChatList)
		case <-components.ChatsRefreshedChan:
			app.QueueUpdateDraw(updateChatList)
		case <-sigChan:
			return
		}
	}
}

func getChats() []models.ChatDTO {
	chats, err := funcs.GetChatsAndInfoWithMessages()
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

func displayMessages(chat models.ChatDTO, chatTextView *tview.TextView) {
	messages := chat.Messages
	ownUser, _, err := funcs.GetOwnUser()
	if err != nil {
		panic(err)
	}
	for _, message := range messages {
		if message.SentByID == ownUser.ID {
			chatTextView.Write(
				[]byte(fmt.Sprintf("\t%s: %s\n\t%s", "You", message.Message, message.SentAt.Format("02 Jan 2006 15:04")+"\n\n")))
		} else {
			chatTextView.Write(
				[]byte(fmt.Sprintf("%s: %s\n%s", chat.ContactName, message.Message, message.SentAt.Format("02 Jan 2006 15:04")+"\n\n")))
		}
	}
}

func BuildMainPage(app *tview.Application, sigChan chan os.Signal) tview.Primitive {
	flex := tview.NewFlex()
	flex.SetBorder(false).SetBorderAttributes(tcell.AttrDim).SetBorderColor(tcell.NewRGBColor(228, 179, 99))

	flex1 := tview.NewFlex()

	upperInfo := components.InfoText()
	logo := components.Logo(tcell.NewRGBColor(232, 233, 235))

	for _, chat := range chats {
		chatListView.AddItem(chat.ContactName, chat.ContactID, 0, nil)
	}

	flex1.SetBorder(true).SetBorderColor(tcell.NewRGBColor(232, 233, 235))
	flex1.SetDirection(tview.FlexRow).AddItem(upperInfo, 0, 1, true)
	flex1.SetDirection(tview.FlexRow).AddItem(logo, 0, 4, true)
	flex1.SetDirection(tview.FlexRow).AddItem(components.LowerTextView, 0, 1, true)

	// Select chat
	chatTextView := tview.NewTextView()
	chatListView.SetSelectedFunc(
		func(row int, mainText, secondaryText string, shortcut rune) {
			chat := getSelectedChat(chats[row].ID)
			chatTitleView := components.ChatTitle(chat.ContactName)
			flex1.Clear()
			flex1.AddItem(chatTitleView, 0, 1, true)
			chatTextView = tview.NewTextView()
			displayMessages(*chat, chatTextView)
			flex1.AddItem(chatTextView, 0, 7, true)
		},
	)

	// Option menu
	optionMenu := components.OptionsMenu(sigChan, flex1)

	flex.SetDirection(tview.FlexColumn).AddItem(chatListView, 0, 1, true)
	flex.SetDirection(tview.FlexColumn).AddItem(flex1, 0, 4, true)
	flex.SetDirection(tview.FlexColumn).AddItem(optionMenu, 0, 1, true)

	go waitForChats(sigChan, app)

	return flex
}
