package pages

import (
	"string_um/string/main/funcs"
	"string_um/string/models"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

var chats = []models.ChatDTO{
	{
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	},
	{
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	}, {
		ID:          uuid.New(),
		ContactID:   "qkdwwd",
		ContactName: "John Doe",
		Messages:    nil,
	},
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

func displaySelectedChat(chatID uuid.UUID) {
	for _, chat := range chats {
		if chat.ID == chatID {
			chatMessages := chat.Messages
			for _, message := range chatMessages {
				lowerTextView.Write([]byte(chat.ContactName + ": " + message.Message + "\n"))
			}
			return
		}
	}
}

func BuildMainPage() tview.Primitive {
	flex := tview.NewFlex()
	flex.SetBorder(false).SetBorderAttributes(tcell.AttrDim).SetBorderColor(tcell.NewRGBColor(228, 179, 99))

	flex1 := tview.NewFlex()

	upperInfo := tview.NewTextView()
	upperInfo.SetText("String - Secure Messaging.\tVersion 1.0.0")
	upperInfo.SetTextStyle(tcell.StyleDefault.Italic(true))
	upperInfo.SetTextAlign(tview.AlignCenter)
	upperInfo.SetTextColor(tcell.NewRGBColor(232, 233, 235))

	logo := tview.NewTextView()
	logo.Write(loadTitle())
	logo.SetTextAlign(tview.AlignCenter)
	logo.SetTextColor(tcell.NewRGBColor(224, 223, 213))

	chatList := tview.NewList()
	chatList.ShowSecondaryText(true)
	chatList.SetMainTextColor(tcell.NewRGBColor(232, 233, 235))
	chatList.SetSecondaryTextColor(tcell.NewRGBColor(241, 217, 177))
	chatList.SetSelectedTextColor(tcell.ColorBlack)
	chatList.SetSelectedBackgroundColor(tcell.NewRGBColor(228, 179, 99))
	chatList.SetBorder(true).SetBorderColor(tcell.NewRGBColor(182, 143, 79))
	chatList.SetTitle("Chats").SetTitleColor(tcell.NewRGBColor(228, 179, 99))

	for _, chat := range chats {
		chatList.AddItem(chat.ContactName, chat.ContactID, 0, nil)
	}

	form := tview.NewForm().SetButtonsAlign(tview.AlignCenter)
	form.SetTitleAlign(tview.AlignCenter)
	form.AddPasswordField("Password: ", "", 30, '*', func(text string) {
		lowerTextView.SetText("")
		inputedPassword = text
	})
	form.AddButton("Login", login)
	form.SetLabelColor(tcell.NewRGBColor(228, 179, 99)).SetFieldBackgroundColor(tcell.NewRGBColor(224, 223, 213)).SetFieldTextColor(tcell.NewRGBColor(50, 55, 57)).SetButtonBackgroundColor(tcell.NewRGBColor(228, 179, 99)).SetButtonTextColor(tcell.ColorBlack)

	flex1.SetBorder(true).SetBorderColor(tcell.NewRGBColor(224, 223, 213))
	flex1.SetDirection(tview.FlexRow).AddItem(upperInfo, 0, 1, true)
	flex1.SetDirection(tview.FlexRow).AddItem(logo, 0, 4, true)

	//flex2.SetDirection(tview.FlexRow).AddItem(logo, 0, 1, true)

	flex.SetDirection(tview.FlexColumn).AddItem(chatList, 0, 1, true)
	flex.SetDirection(tview.FlexColumn).AddItem(flex1, 0, 4, true)
	flex.SetDirection(tview.FlexColumn).AddItem(tview.NewBox(), 0, 1, true)

	return flex
}
