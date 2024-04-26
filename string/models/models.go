package models

import (
	"time"

	"github.com/google/uuid"
)

type OwnUser struct {
	ID         string `gorm:"type:string;primaryKey;unique;not null;serializer:json" json:"id"`
	Password   string `gorm:"type:string;not null;serializer:json" json:"password"`
	PrivateKey []byte `gorm:"type:byte[];not null;serializer:json" json:"privateKey"`
}

type Contact struct {
	ID   string `gorm:"type:string;primaryKey;unique;not null;serializer:json" json:"id"`
	Name string `gorm:"type:string;not null;serializer:json" json:"name"`
}

type ContactAddress struct {
	ID              uuid.UUID `gorm:"type:string;primaryKey;unique;not null;serializer:json" json:"id"`
	Contact         Contact   `gorm:"type:embedded;not null;serializer:json" json:"contact"`
	ObservedAddress string    `gorm:"type:string;serializer:json" json:"observedAddress"`
	ObservedAt      time.Time `gorm:"type:timestamp;serializer:json" json:"observedAt"`
}

type Chat struct {
	ID      uuid.UUID `gorm:"type:string;primaryKey;unique;not null;serializer:json" json:"id"`
	Contact Contact   `gorm:"type:embedded;not null;serializer:json" json:"contact"`
}

type Message struct {
	ID         uuid.UUID `gorm:"type:string;primaryKey;unique;not null;serializer:json" json:"id"`
	Chat       Chat      `gorm:"type:embedded;not null;serializer:json" json:"chat"`
	AlredySent bool      `gorm:"type:boolean;not null;serializer:json" json:"alredySent"`
	SentBy     Contact   `gorm:"type:embedded;serializer:json" json:"sentBy"`
	SentAt     time.Time `gorm:"type:timestamp;serializer:json" json:"sentAt"`
	Message    string    `gorm:"type:string;serializer:json" json:"message"`
}

type MessageDTO struct {
	ID      uuid.UUID `gorm:"type:string;primaryKey;unique;not null;serializer:json" json:"id"`
	Chat    Chat      `gorm:"type:embedded;not null;serializer:json" json:"chat"`
	SentBy  Contact   `gorm:"type:embedded;serializer:json" json:"sentBy"`
	SentAt  time.Time `gorm:"type:timestamp;serializer:json" json:"sentAt"`
	Message string    `gorm:"type:string;serializer:json" json:"message"`
}
