package components

import (
	"string_um/string/networking/node"

	"github.com/gdamore/tcell/v2"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/rivo/tview"
)

var addContactForm = tview.NewForm()
var libp2pHost host.Host

func addContact() {
	name := addContactForm.GetFormItemByLabel("Name: ").(*tview.InputField).GetText()
	multihash := addContactForm.GetFormItemByLabel("Multihash: ").(*tview.InputField).GetText()
	if name == "" {
		LowerTextView.SetText("Name can't be empty.")
		return
	} else if multihash == "" {
		LowerTextView.SetText("Multihash can't be empty.")
		return

	}
	if err := node.AddNewContact(libp2pHost, multihash, name); err != nil {
		LowerTextView.SetText(err.Error())
	} else {
		LowerTextView.SetText("Contact added.")
	}
}

func goBack() {
	containingFlex.Clear()
	containingFlex.AddItem(InfoText(), 0, 1, true)
	containingFlex.AddItem(Logo(tcell.NewRGBColor(232, 233, 235)), 0, 4, true)
}

func AddContactForm(p2phost host.Host) *tview.Form {
	libp2pHost = p2phost

	LowerTextView.SetText("")
	addContactForm.SetBorder(true)
	addContactForm.SetTitle("Add new contact")
	addContactForm.SetTitleAlign(tview.AlignCenter)
	addContactForm.SetTitleColor(tcell.NewRGBColor(232, 233, 235))
	addContactForm.SetBorderColor(tcell.NewRGBColor(114, 89, 49))
	addContactForm.SetBackgroundColor(tcell.ColorBlack)
	addContactForm.SetLabelColor(tcell.NewRGBColor(232, 233, 235))
	addContactForm.SetFieldBackgroundColor(tcell.NewRGBColor(224, 223, 213))
	addContactForm.SetFieldTextColor(tcell.NewRGBColor(50, 55, 57))
	addContactForm.SetButtonBackgroundColor(tcell.NewRGBColor(114, 89, 49))
	addContactForm.SetButtonTextColor(tcell.NewRGBColor(224, 223, 213))
	addContactForm.SetButtonsAlign(tview.AlignCenter)
	addContactForm.SetCancelFunc(goBack)
	addContactForm.AddInputField("Name: ", "", 40, nil, nil)
	addContactForm.AddInputField("Multihash: ", "", 40, nil, nil)
	addContactForm.AddButton("Add", addContact)
	addContactForm.AddFormItem(LowerTextView)
	return addContactForm
}
