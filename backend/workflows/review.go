package workflows

import (
	"backend/config"
	"backend/db"
	"backend/db/models"
	"backend/nlp"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

func reviewReminder(body map[string]string) string {
	ConversationID, _ := strconv.ParseUint(body["conversationId"], 10, 32)
	var conversation models.Conversation
	db.DB.First(&conversation, body["conversationId"])

	var contextMap map[string]interface{}
	if err := json.Unmarshal(conversation.Context, &contextMap); err != nil {
		log.Fatalf("Failed to unmarshal Context: %v", err)
	}

	type Context struct {
		ProductID string `json:"key"`
	}

	var contextStruct Context
	if err := json.Unmarshal(conversation.Context, &contextStruct); err != nil {
		log.Fatalf("Failed to unmarshal Context into struct: %v", err)
	}

	var product models.Product
	db.DB.First(&product, contextStruct.ProductID)

	var message = fmt.Sprintf("Hello again! We noticed you've recently received your %s. We'd love to hear about your experience. Can you spare a few minutes to share your thoughts?", product.Name)

	newSystemMessage := models.Message{
		Content:        message,
		IsUserInput:    false,
		UserID:         uint(config.BOT_USER_ID),
		ConversationID: uint(ConversationID),
	}
	db.DB.Create(&newSystemMessage)

	conversation.WorkflowStatus = models.PRODUCT_REVIEW_CONFIRMATION
	db.DB.Save(&conversation)

	return message
}

func reviewConfirmation(body map[string]string) string {
	ConversationID, _ := strconv.ParseUint(body["conversationId"], 10, 32)
	var conversation models.Conversation
	db.DB.First(&conversation, body["conversationId"])

	var contextMap map[string]interface{}
	if err := json.Unmarshal(conversation.Context, &contextMap); err != nil {
		log.Fatalf("Failed to unmarshal Context: %v", err)
	}

	type Context struct {
		ProductID string `json:"key"`
	}

	var contextStruct Context
	if err := json.Unmarshal(conversation.Context, &contextStruct); err != nil {
		log.Fatalf("Failed to unmarshal Context into struct: %v", err)
	}

	var product models.Product
	db.DB.First(&product, contextStruct.ProductID)

	affirmative := nlp.IsMessageAffirmative(body["message"])

	if affirmative {
		message := fmt.Sprintf("Fantastic! On a scale of 1-5, how would you rate the %s?", product.Name)
		newSystemMessage := models.Message{
			Content:        message,
			IsUserInput:    false,
			UserID:         uint(config.BOT_USER_ID),
			ConversationID: uint(ConversationID),
		}
		db.DB.Create(&newSystemMessage)

		conversation.WorkflowStatus = models.PRODUCT_REVIEW_RATING
		db.DB.Save(&conversation)

		return message
	} else {
		message := fmt.Sprintf("Not a problem, enjoy your %s!", product.Name)
		newSystemMessage := models.Message{
			Content:        message,
			IsUserInput:    false,
			UserID:         uint(config.BOT_USER_ID),
			ConversationID: uint(ConversationID),
		}
		db.DB.Create(&newSystemMessage)

		conversation.WorkflowStatus = models.PRODUCT_REVIEW_FINISHED
		db.DB.Save(&conversation)

		return message
	}
}

func reviewRating(body map[string]string) string {
	ConversationID, _ := strconv.ParseUint(body["conversationId"], 10, 32)
	var conversation models.Conversation
	db.DB.First(&conversation, body["conversationId"])

	hasRating, Rating, _ := nlp.AnalyzeUserRating(body["message"])

	if hasRating {
		message := "Thank you for sharing your feedback! If you have any more thoughts or need assistance with anything else, feel free to reach out!"
		newSystemMessage := models.Message{
			Content:        message,
			IsUserInput:    false,
			UserID:         uint(config.BOT_USER_ID),
			ConversationID: uint(ConversationID),
		}
		db.DB.Create(&newSystemMessage)

		var contextMap map[string]interface{}
		if err := json.Unmarshal(conversation.Context, &contextMap); err != nil {
			log.Fatalf("Failed to unmarshal Context: %v", err)
		}

		type Context struct {
			ProductID int `json:"ProductID"`
		}

		var contextStruct Context
		if err := json.Unmarshal(conversation.Context, &contextStruct); err != nil {
			log.Fatalf("Failed to unmarshal Context into struct: %v", err)
		}

		newProductReview := models.ProductReview{
			UserID:    conversation.UserID,
			ProductID: uint(contextStruct.ProductID),
			Rating:    uint(Rating),
		}
		db.DB.Create(&newProductReview)

		conversation.WorkflowStatus = models.PRODUCT_REVIEW_FINISHED
		db.DB.Save(&conversation)
		return message
	} else {
		message := "Sorry I didn't understand your message, could you please share your rating again?"
		newSystemMessage := models.Message{
			Content:        message,
			IsUserInput:    false,
			UserID:         uint(config.BOT_USER_ID),
			ConversationID: uint(ConversationID),
		}
		db.DB.Create(&newSystemMessage)
		return message
	}
}
