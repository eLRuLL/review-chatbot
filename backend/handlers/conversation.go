package handlers

import (
	"backend/db"
	"backend/db/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func CreateConversation(c *gin.Context) {
	var newConversation models.Conversation

	if err := c.ShouldBindJSON(&newConversation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contextMap := map[string]interface{}{
		"ProductID": 1,
	}
	contextJSON, _ := json.Marshal(contextMap)

	newConversation.WorkflowStatus = models.PRODUCT_REVIEW
	newConversation.Context = datatypes.JSON(contextJSON)

	db.DB.Create(&newConversation)
	c.JSON(http.StatusOK, gin.H{"data": newConversation})
}

func GetConversation(c *gin.Context) {
	var conversation models.Conversation
	id := c.Param("id")

	if result := db.DB.First(&conversation, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": conversation})
}
