package db

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"strconv"
)

type MessageRepository interface {
	Save(ctx context.Context, message *Message) error
	Saves(ctx context.Context, message []Message) error
	Read(ctx context.Context, messageID int) (*Message, error)
}

type fakeMessage struct {
}

func NewMessageRepository() MessageRepository {
	return &fakeMessage{}
}

func (o *fakeMessage) Save(ctx context.Context, message *Message) error {
	log.Printf("Saved message: %s", message.ID)
	return nil
}

func (o *fakeMessage) Saves(ctx context.Context, message []Message) error {
	log.Printf("Saved messages: %d", len(message))
	return nil
}

func (o *fakeMessage) Read(ctx context.Context, messageID int) (*Message, error) {
	log.Printf("Read message: %d", messageID)
	return &Message{}, nil
}

type MessageTemplateRepository interface {
	QueryOne(ctx context.Context, sentimentType MessageTemplateSentimentType) (*MessageTemplate, error)
}

type fakeMessageTemplate struct {
}

func NewMessageTemplateRepository() MessageTemplateRepository {
	return &fakeMessageTemplate{}
}

func (o *fakeMessageTemplate) QueryOne(ctx context.Context, sentimentType MessageTemplateSentimentType) (*MessageTemplate, error) {
	messageText := ""

	switch sentimentType {
	case MessageTemplateSentimentTypeNegative:
		messageText = NegativeMessages[rand.Intn(len(NegativeMessages))]
	case MessageTemplateSentimentTypePositive:
		messageText = PositiveMessages[rand.Intn(len(PositiveMessages))]
	default:
		return nil, errors.New("invalid sentiment type")
	}

	result := &MessageTemplate{
		ID:   "TemplateID" + strconv.Itoa(rand.Intn(10000)),
		Text: messageText,
		Type: MessageTemplateTypeText,

		SentimentType: sentimentType,
	}

	return result, nil
}

var PositiveMessages = []string{
	"Thank you for choosing us! Your order has been successfully placed.",
	"Congratulations on your purchase! Your items are on their way to you.",
	"Your shopping experience doesn't end here! We're thrilled to have you as our customer.",
}

var NegativeMessages = []string{
	"We're truly sorry for the inconvenience. Rest assured, we're taking steps to make things right",
	"We deeply regret any inconvenience this has caused you. We're actively investigating and will update you soon.",
	"Please accept our apologies for the inconvenience. We're committed to learning from this and improving our services.",
}
