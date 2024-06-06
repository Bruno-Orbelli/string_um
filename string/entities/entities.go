package entities

import (
	"time"

	"github.com/google/uuid"
)

type OwnUser struct {
	ID           string `gorm:"type:string;primaryKey;unique;not null" json:"id"`
	PasswordHash string `gorm:"type:string;not null" json:"passwordHash"`
	EncodingHash string `gorm:"type:string;not null" json:"encodingHash"`
	PrivateKey   []byte `gorm:"type:byte[];not null" json:"privateKey"`
}

type Contact struct {
	ID               string           `gorm:"type:string;primaryKey;unique;not null" json:"id"`
	Name             string           `gorm:"type:string;not null" json:"name"`
	ContactAddresses []ContactAddress `gorm:"foreignKey:ContactID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"contactAddresses"`
	Chat             Chat             `gorm:"foreignKey:ContactID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"chat"`
}

type ContactAddress struct {
	ID              uuid.UUID `gorm:"type:string;primaryKey;unique;not null" json:"id"`
	ContactID       string    `gorm:"not null;constraint:OnDelete:CASCADE" json:"contactID"`
	ObservedAddress string    `gorm:"type:string" json:"observedAddress"`
	ObservedAt      time.Time `gorm:"type:timestamp" json:"observedAt"`
}

type Chat struct {
	ID        uuid.UUID `gorm:"type:string;primaryKey;unique;not null" json:"id"`
	ContactID string    `gorm:"type:string;not null;constraint:OnDelete:CASCADE" json:"contactID"`
	Messages  []Message `gorm:"foreignKey:ChatID" json:"messages"`
}

type ChatDTO struct {
	ID          uuid.UUID `gorm:"type:string;primaryKey;unique;not null" json:"id"`
	ContactID   string    `gorm:"type:string;not null" json:"contactID"`
	ContactName string    `gorm:"type:string;not null" json:"contactName"`
	Messages    []Message `gorm:"foreignKey:ChatID" json:"messages"`
}

type Message struct {
	ID          uuid.UUID `gorm:"type:string;primaryKey;unique;not null" json:"id"`
	ChatID      uuid.UUID `gorm:"type:string;not null;constraint:OnDelete:CASCADE" json:"chatID"`
	AlreadySent bool      `gorm:"type:boolean;not null" json:"alreadySent"`
	SentByID    string    `gorm:"type:string;not null;constraint:OnDelete:CASCADE" json:"sentByID"`
	SentAt      time.Time `gorm:"type:timestamp" json:"sentAt"`
	Message     string    `gorm:"type:string" json:"message"`
}

type MessageDTO struct {
	ID       uuid.UUID `gorm:"type:string;primaryKey;unique;not null" json:"id"`
	ChatID   uuid.UUID `gorm:"not null;constraint:OnDelete:CASCADE" json:"chat"`
	SentByID string    `gorm:"constraint:OnDelete:CASCADE;serializer:json" json:"sentBy"`
	SentAt   time.Time `gorm:"type:timestamp" json:"sentAt"`
	Message  string    `gorm:"type:string" json:"message"`
}
