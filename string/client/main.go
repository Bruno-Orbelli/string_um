package main

import (
	"fmt"
	dapi "string_um/string/client/debug-api"
)

func main() {
	fmt.Println("Starting database API...")
	dapi.RunDatabaseAPI()
}
