package prod_api

import (
	"fmt"
	"string_um/string/entities"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Database *gorm.DB

// RunDatabaseAPI initializes the database and runs the API
func RunDatabaseAPI() {
	var err error
	Database, err = gorm.Open(
		sqlite.Open("test.db"),
		&gorm.Config{TranslateError: true},
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// Auto-migrate the database
	if err = Database.AutoMigrate(
		&entities.Chat{},
		&entities.Contact{},
		&entities.ContactAddress{},
		&entities.Message{},
		&entities.OwnUser{},
	); err != nil {
		panic(fmt.Sprintf("Failed to auto-migrate database: %v", err))
	}

	// fmt.Println("Database direct API running.")
}
