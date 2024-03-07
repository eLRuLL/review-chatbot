package validation

type MessageCreateSchema struct {
	Content        string `json:"content" binding:"required"`
	ConversationID uint   `json:"conversationId" binding:"required"`
}
