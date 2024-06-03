package main

import (
	"context"
	"os"
	"os/signal"
	prod_api "string_um/string/client/prod-api"
	"string_um/string/main/funcs"
	"string_um/string/main/tui/components"
	"string_um/string/main/tui/globals"
	pag "string_um/string/main/tui/pages"
	"string_um/string/models"
	"syscall"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
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

func startHost(ctx context.Context, ownUser models.OwnUser) (host.Host, error) {
	privKey, err := crypto.UnmarshalPrivateKey(ownUser.PrivateKey)
	if err != nil {
		return nil, err
	}
	host, _, err := funcs.StartHost(ctx, privKey)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = network.WithAllowLimitedConn(ctx, "relay info")
	defer cancel()

	prod_api.RunDatabaseAPI()
	app := tview.NewApplication()
	boolForReg := isRegistered()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	globals.Pages.AddPage("register", pag.BuildRegisterPage(), true, boolForReg[0])
	globals.Pages.AddPage("login", pag.BuildLoginPage(), true, boolForReg[1])
	globals.Pages.AddPage("main", pag.BuildMainPage(app, sigChan), true, false)

	go func() { // Wait for login to be successful and get private key for node setup
		<-globals.LoginSuccessChan
		ownUser, _, err := funcs.GetOwnUser()
		if err != nil {
			panic(err)
		}
		host, err := startHost(ctx, *ownUser)
		if err != nil {
			panic(err)
		}
		components.Libp2pHost = host
		globals.ChatsReadyChan <- true
	}()

	go func() {
		<-sigChan
		close(sigChan) // Signal the waitForChats goroutine to stop
		cancel()
		<-ctx.Done()
		app.Stop()
	}()

	defer func() {
		if err := funcs.CloseDatabase(); err != nil {
			panic(err)
		}
	}()

	if err := app.SetRoot(globals.Pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
