package workflows

import (
	"backend/db"
	"backend/db/models"
)

func ReplyHandler(body map[string]string) string {
	var conversation models.Conversation
	db.DB.First(&conversation, body["conversationId"])

	if conversation.WorkflowStatus == models.PRODUCT_REVIEW {
		return reviewReminder(body)
	}

	if conversation.WorkflowStatus == models.PRODUCT_REVIEW_CONFIRMATION {
		return reviewConfirmation(body)
	}

	if conversation.WorkflowStatus == models.PRODUCT_REVIEW_RATING {
		return reviewRating(body)
	}

	return "Internal Error"
}
