package components

import (
	"os"
	"path/filepath"

	"github.com/rivo/tview"
)

var LowerTextView = tview.NewTextView().SetTextAlign(tview.AlignCenter)
var Pages = tview.NewPages()

var ChatsReadyChan = make(chan bool, 1)
var ChatsRefreshedChan = make(chan bool, 1)

func LoadTitle() []byte {
	path := filepath.Join(".", "resources", "string.txt")
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return b
}
