package model

import "messageBot/db"

type Message struct {
	Object string         `json:"object"`
	Entry  []MessageEntry `json:"entry"`

	SentimentType db.MessageTemplateSentimentType `json:"sentiment_typ,omitempty"`
}

type MessageEntry struct {
	ID          string        `json:"id"`
	Time        int           `json:"time"`
	MessageData []MessageData `json:"messaging"`
}

type MessageData struct {
	Sender    MessageDataSender    `json:"sender"`
	Recipient MessageDataRecipient `json:"recipient"`
	Message   MessageDataMessage   `json:"message,omitempty"`
}

func (m *Message) GetRecipientID() string {
	if len(m.Entry) > 0 && len(m.Entry[0].MessageData) > 0 {
		return m.Entry[0].MessageData[0].Recipient.ID
	}
	return ""
}

func (m *Message) GetText() string {
	result := ""
	for _, e := range m.Entry {
		for _, data := range e.MessageData {
			result += data.Message.Text
		}
	}
	return result
}

type MessageDataMessage struct {
	Mid  string `json:"mid"`
	Seq  int    `json:"seq"`
	Text string `json:"text"`
}

type MessageDataSender struct {
	ID string `json:"id"`
}

type MessageDataRecipient struct {
	ID string `json:"id"`
}

const (
	SendAPIEndpoint = "https://graph.facebook.com/v2.6/me/messages"
)

var (
	// Prod if run in prod  model
	// will not allow to connect fb or huggingface.
	Prod bool

	VerifyToken     = "verify_token"
	PageAccessToken = "page_access_token"
	AppSecret       = "app_secret"
	NLPToken        = "nlp_token"
)

type MessageResponse struct {
	MessageType string               `json:"message_type"`
	Recipient   MessageDataRecipient `json:"recipient"`
}
