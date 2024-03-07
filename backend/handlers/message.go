package handlers

import (
	"backend/db"
	"backend/db/models"
	"backend/workflows"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateMessage(c *gin.Context) {
	var newUserMessage models.Message

	conversationID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := c.ShouldBindJSON(&newUserMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUserMessage.IsUserInput = true
	newUserMessage.ConversationID = uint(conversationID)
	db.DB.Create(&newUserMessage)

	replyBody := map[string]string{
		"conversationId": c.Param("id"),
		"message":        newUserMessage.Content,
	}

	newSystemMessage := workflows.ReplyHandler(replyBody)

	c.JSON(http.StatusOK, gin.H{"reply": newSystemMessage})
}
