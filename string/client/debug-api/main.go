package debug_api

import (
	"fmt"
	"net/http"
	"string_um/string/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Database *gorm.DB

func getMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Contacts
	mux.HandleFunc("/contacts", GetContacts)
	mux.HandleFunc("/contacts/{id}", GetContact)
	mux.HandleFunc("/contacts/create", CreateContact)
	mux.HandleFunc("/contacts/delete/{id}", DeleteContact)

	// Contact Addresses
	mux.HandleFunc("/contactAddresses", GetContactAddresses)
	mux.HandleFunc("/contactAddresses/{id}", GetContactAddress)
	mux.HandleFunc("/contactAddresses/create", CreateContactAddress)
	mux.HandleFunc("/contactAddresses/delete/{id}", DeleteContactAddress)

	// Chats
	mux.HandleFunc("/chats", GetChats)
	mux.HandleFunc("/chats/{id}", GetChat)
	mux.HandleFunc("/chats/create", CreateChat)
	mux.HandleFunc("/chats/delete/{id}", DeleteChat)

	// Messages
	mux.HandleFunc("/messages", GetMessages)
	mux.HandleFunc("/messages/{id}", GetMessage)
	mux.HandleFunc("/messages/create", CreateMessage)
	mux.HandleFunc("/messages/update/{id}", UpdateMessage)
	mux.HandleFunc("/messages/delete/{id}", DeleteMessage)

	// Own User
	mux.HandleFunc("/ownUser", GetOwnUser)
	mux.HandleFunc("/ownUser/create", CreateOwnUser)
	mux.HandleFunc("/ownUser/delete", DeleteOwnUser)

	return mux
}

func RunDatabaseAPI() {
	var err error
	mux := getMux()
	Database, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{TranslateError: true})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// Auto-migrate the database
	if err = Database.AutoMigrate(
		&models.Chat{},
		&models.Contact{},
		&models.ContactAddress{},
		&models.Message{},
		&models.OwnUser{},
	); err != nil {
		panic(fmt.Sprintf("Failed to auto-migrate database: %v", err))
	}

	fmt.Println("Database API running on port 3000")
	http.ListenAndServe(":3000", mux)
}
