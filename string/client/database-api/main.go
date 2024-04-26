package database_api

import (
	"fmt"
	"net/http"
)

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
	mux.HandleFunc("/messages/setAsSent/{id}", UpdateMessageToSent)
	mux.HandleFunc("/messages/delete/{id}", DeleteMessage)

	// Own User
	mux.HandleFunc("/ownUser", GetOwnUser)
	mux.HandleFunc("/ownUser/create", CreateOwnUser)
	mux.HandleFunc("/ownUser/delete", DeleteOwnUser)

	return mux
}

func RunDatabaseAPI() {
	mux := getMux()
	http.ListenAndServe(":3000", mux)
	fmt.Println("Database API running on port 3000")
}
