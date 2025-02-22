package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func AddedContactsList() *tview.List {
	addedContacts := tview.NewList()
	addedContacts.ShowSecondaryText(true)
	addedContacts.SetTitle("Added contacts")
	addedContacts.SetTitleAlign(tview.AlignCenter)
	addedContacts.SetTitleColor(tcell.NewRGBColor(232, 233, 235))
	addedContacts.SetSecondaryTextColor(tcell.NewRGBColor(111, 115, 116))
	addedContacts.SetMainTextColor(tcell.NewRGBColor(232, 233, 235))
	addedContacts.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	addedContacts.SetSelectedTextColor(tcell.NewRGBColor(232, 233, 235))
	addedContacts.SetSelectedBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	addedContacts.SetBorder(true)
	addedContacts.SetBorderColor(tcell.NewRGBColor(51, 255, 106))

	return addedContacts
}
