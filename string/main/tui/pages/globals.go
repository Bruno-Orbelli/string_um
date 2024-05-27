package pages

import (
	"os"
	"path/filepath"

	"github.com/rivo/tview"
)

var inputedPassword = ""
var lowerTextView = tview.NewTextView().SetTextAlign(tview.AlignCenter)
var Pages = tview.NewPages()

func loadTitle() []byte {
	path := filepath.Join(".", "resources", "string.txt")
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return b
}
