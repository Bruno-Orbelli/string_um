package main

import (
	"fmt"
	dapi "string_um/string/client/database-api"
)

func main() {
	fmt.Print("Starting database API...")
	dapi.RunDatabaseAPI()
}
