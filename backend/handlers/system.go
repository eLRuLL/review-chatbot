package handlers

import (
	"backend/config"
	"backend/db"
	"backend/db/models"
)

func createSystemMessage(conversationID uint) models.Message {
	systemMessageContent := "This is an automated response from the system."
	newSystemMessage := models.Message{
		Content:        systemMessageContent,
		IsUserInput:    false,
		UserID:         uint(config.BOT_USER_ID),
		ConversationID: conversationID,
	}

	db.DB.Create(&newSystemMessage)
	return newSystemMessage
}
