package main

import (
	"os"
	pag "string_um/string/main/tui/pages"

	"github.com/rivo/tview"
)

func isRegistered() []bool {
	boolArr := make([]bool, 2)
	if _, err := os.Stat("en_test.db"); err != nil {
		if os.IsNotExist(err) {
			boolArr[0] = true
			boolArr[1] = false
		} else {
			panic(err)
		}
	} else {
		boolArr[0] = false
		boolArr[1] = true
	}
	return boolArr
}

func main() {
	app := tview.NewApplication()
	boolForReg := isRegistered()

	pag.Pages.AddPage("register", pag.BuildRegisterPage(), true, boolForReg[0])
	pag.Pages.AddPage("login", pag.BuildLoginPage(), true, boolForReg[1])
	pag.Pages.AddPage("main", pag.BuildMainPage(), true, false)

	if err := app.SetRoot(pag.Pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
