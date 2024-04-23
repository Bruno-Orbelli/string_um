package models

import (
	"time"

	"github.com/google/uuid"
)

type OwnUser struct {
	ID         string `gorm:"type:string;primaryKey;unique;not null;serializer:json"`
	Username   string `gorm:"type:string;not null;serializer:json"`
	Password   string `gorm:"type:string;not null;serializer:json"`
	PrivateKey string `gorm:"type:string;not null;serializer:json"`
}

type Contact struct {
	ID       string `gorm:"type:string;primaryKey;unique;not null;serializer:json"`
	Username string `gorm:"type:string;not null;serializer:json"`
}

type ContactAddress struct {
	ID                   uuid.UUID `gorm:"type:string;primaryKey;unique;not null;serializer:json"`
	Contact              Contact   `gorm:"type:embedded;not null;serializer:json"`
	ObservedMultiaddress string    `gorm:"type:string;serializer:json"`
	ObservedAt           time.Time `gorm:"type:timestamp;serializer:json"`
}

type Chat struct {
	ID      uuid.UUID `gorm:"type:string;primaryKey;unique;not null;serializer:json"`
	Contact Contact   `gorm:"type:embedded;not null;serializer:json"`
}

type Message struct {
	ID         uuid.UUID `gorm:"type:string;primaryKey;unique;not null;serializer:json"`
	Chat       Chat      `gorm:"type:embedded;not null;serializer:json"`
	AlredySent bool      `gorm:"type:boolean;not null;serializer:json"`
	SentBy     Contact   `gorm:"type:embedded;serializer:json"`
	SentAt     time.Time `gorm:"type:timestamp;serializer:json"`
	Message    string    `gorm:"type:string;serializer:json"`
}

type MessageDTO struct {
	ID      uuid.UUID `gorm:"type:string;primaryKey;unique;not null;serializer:json"`
	Chat    Chat      `gorm:"type:embedded;not null;serializer:json"`
	SentBy  Contact   `gorm:"type:embedded;serializer:json"`
	SentAt  time.Time `gorm:"type:timestamp;serializer:json"`
	Message string    `gorm:"type:string;serializer:json"`
}
