package components

import (
	"fmt"
)

func Message(isOwn bool, senderName string, body string) string {
	var message string
	if isOwn {
		message = fmt.Sprintf("[green]%s: %s", senderName, body)
	} else {
		message = fmt.Sprintf("[#e8e9eb]%s: %s", senderName, body)
	}

	return message
}
