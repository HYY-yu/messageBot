package db

// Message the message table of database, used to save the message which is sent by fb webhook.
type Message struct {
	ID          string `json:"id" gorm:"primary_key;column:id"`
	SenderID    string `json:"sender_id" gorm:"column:sender_id"`
	RecipientID string `json:"recipient_id" gorm:"column:recipient_id"`
	Timestamp   int    `json:"timestamp" gorm:"column:timestamp"`

	MessageJson string `json:"message_json" gorm:"column:message_json"`
}

// MessageTemplate the message template table of database, a bot can use a template to response a Message.
type MessageTemplate struct {
	ID             string                        `json:"id" gorm:"primary_key;column:id"`
	Type           MessageTemplateType           `json:"type" gorm:"column:type"`
	Text           string                        `json:"text" gorm:"column:text"`
	AttachmentType MessageTemplateAttachmentType `json:"attachment_type" gorm:"column:attachment_type"`
	// AttachmentPayloadJson the fb message payload
	// see on: https://developers.facebook.com/docs/messenger-platform/reference/templates/generic
	AttachmentPayloadJson string `json:"attachment_payload_json" gorm:"column:attachment_payload_json"`
	IsAdvertisement       bool   `json:"is_advertisement" gorm:"column:is_advertisement"` // an advertisement message mark the template contain promotional content

	SentimentType MessageTemplateSentimentType `json:"sentiment_type" gorm:"column:sentiment_type"`
}

type MessageTemplateType int

const (
	MessageTemplateTypeText MessageTemplateType = iota + 1
	MessageTemplateTypeAttachment
)

type MessageTemplateAttachmentType string

const (
	MessageTemplateAttachmentTypeImage    MessageTemplateAttachmentType = "image"
	MessageTemplateAttachmentTypeVideo                                  = "video"
	MessageTemplateAttachmentTypeAudio                                  = "audio"
	MessageTemplateAttachmentTypeFile                                   = "file"
	MessageTemplateAttachmentTypeTemplate                               = "template"
)

type MessageTemplateSentimentType string

const (
	MessageTemplateSentimentTypePositive MessageTemplateSentimentType = "positive"
	MessageTemplateSentimentTypeNegative                              = "negative"
)
