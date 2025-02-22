package prod_api

import (
	"fmt"
	"log"
	"os"
	"string_um/string/entities"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

// RunDatabaseAPI initializes the database and runs the API
func RunDatabaseAPI() {
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	Database, err = gorm.Open(
		sqlite.Open("test.db"),
		&gorm.Config{
			TranslateError: true,
			Logger:         newLogger,
		},
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
