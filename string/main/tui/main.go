package main

import (
	"os"
	"os/signal"
	prod_api "string_um/string/client/prod-api"
	"string_um/string/main/funcs"
	"string_um/string/main/tui/components"
	pag "string_um/string/main/tui/pages"
	"syscall"

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
	prod_api.RunDatabaseAPI()
	app := tview.NewApplication()
	boolForReg := isRegistered()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	components.Pages.AddPage("register", pag.BuildRegisterPage(), true, boolForReg[0])
	components.Pages.AddPage("login", pag.BuildLoginPage(), true, boolForReg[1])
	components.Pages.AddPage("main", pag.BuildMainPage(app, sigChan), true, false)

	go func() {
		<-sigChan
		close(sigChan) // Signal the waitForChats goroutine to stop
		app.Stop()
	}()

	defer func() {
		if err := funcs.CloseDatabase(); err != nil {
			panic(err)
		}
	}()

	if err := app.SetRoot(components.Pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
