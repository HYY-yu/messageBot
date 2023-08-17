package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"messageBot/messenger/model"
	"strconv"
)

var DB *gorm.DB

type MessageRepository interface {
	Saves(ctx context.Context, message []Message) error
}

type fakeMessage struct {
}

func NewMessageRepository() MessageRepository {
	return &fakeMessage{}
}

func (o *fakeMessage) Saves(ctx context.Context, message []Message) error {
	if model.Prod {
		err := DB.Create(message).Error
		if err != nil {
			return err
		}

		return nil
	}
	log.Printf("Saved messages: %d", len(message))
	return nil
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
	"Buckle up for a fantastic ride! Your shopping adventure doesn't stop here, and we're thrilled to accompany you.",
	"Your choice has us smiling! Sit tight as your order embarks on its journey to you.",
	"Hooray for your purchase! Your items are en route, and we're excited to have you as part of our family.",
	"A big thank you for selecting us! Your order is all set and set to reach you soon.",
	"But wait, there's more! Your shopping journey continues with us. We're delighted to have you onboard.",
	"Kudos on your new purchase! Get ready to enjoy your items as they make their way to your doorstep.",
	"We're truly grateful you picked us! Your order has been confirmed and is on its way.",
}

var NegativeMessages = []string{
	"We're truly sorry for the inconvenience. Rest assured, we're taking steps to make things right",
	"We deeply regret any inconvenience this has caused you. We're actively investigating and will update you soon.",
	"Please accept our apologies for the inconvenience. We're committed to learning from this and improving our services.",
	"Apologies for any inconvenience caused. We're committed to addressing this promptly and ensuring a smoother experience in the future.",
	"We're genuinely sorry for the disruption this has caused. We're determined to turn this around and regain your trust.",
	"We deeply regret the inconvenience you've experienced. We're fully committed to finding a resolution that meets your expectations.",
	"Please accept our sincerest apologies for the hassle. We're actively seeking ways to prevent this from happening again.",
	"We understand the frustration this may have caused. Rest assured, we're striving to set things right.",
	"Our heartfelt apologies for the inconvenience. We're already working on a solution to rectify this situation.",
	"We sincerely apologize for any trouble，we will do our best to make things right.",
}
