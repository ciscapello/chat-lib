package contracts

import "github.com/google/uuid"

type UserCreatedMessage struct {
	Email    string
	Username string
	Code     string
}

type MessageCreatedBody struct {
	SenderId       uuid.UUID
	ConversationId int
	MessageBody    string
}

type MessageSocketBody struct {
	Type           string `json:"type"`
	ConversationId int    `json:"conversation_id,omitempty"`
	FromUserID     string `json:"from_user_id,omitempty"`
	ToUserID       string `json:"to_user_id,omitempty"`
	MessageBody    string `json:"message_body,omitempty"`
}

const (
	UserCreatedTopic    = "user.created"
	MessageCreatedTopic = "message.created"
)
