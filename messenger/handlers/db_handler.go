package handlers

import (
	"context"
	"messageBot/db"
	"messageBot/db/model"
)

// DatabaseMessageHandler save message to database
type DatabaseMessageHandler struct {
	Repo db.MessageRepository
}

func NewDatabaseMessageHandler(repo db.MessageRepository) *DatabaseMessageHandler {
	return &DatabaseMessageHandler{
		Repo: repo,
	}
}

func (h *DatabaseMessageHandler) Handle(ctx context.Context, message *model.Message) error {
	messageDbs := make([]model.MessageTable, 0)
	for _, e := range message.Entry {
		for _, data := range e.MessageData {
			dbMsg := &model.MessageTable{
				SenderID:    data.Sender.ID,
				RecipientID: data.Recipient.ID,
				Timestamp:   e.Time,
				MessageText: data.Message.Text,
			}
			messageDbs = append(messageDbs, *dbMsg)
		}
	}
	return h.Repo.Saves(ctx, messageDbs)
}
