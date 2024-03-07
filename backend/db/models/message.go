package models

import "gorm.io/gorm"

type Status string

const (
	CREATED  Status = "CREATED"
	QUEUED   Status = "QUEUED"
	SENT     Status = "SENT"
	RECEIVED Status = "RECEIVED"
)

type Message struct {
	gorm.Model

	UserID         uint
	ConversationID uint
	Content        string
	IsUserInput    bool
	Status         Status

	User         User         `gorm:"foreignKey:UserID"`
	Conversation Conversation `gorm:"foreignKey:ConversationID"`
}
