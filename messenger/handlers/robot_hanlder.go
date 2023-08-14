package handlers

import (
	"bytes"
	"context"
	"io"
	"log"
	"messageBot/db"
	"messageBot/messenger"
	"messageBot/messenger/model"
	"net/http"
)

type RobotHandler struct {
	Repo            db.MessageTemplateRepository
	PageAccessToken string
}

func NewRobotHandler(accessToken string, repo db.MessageTemplateRepository) *RobotHandler {
	return &RobotHandler{
		Repo:            repo,
		PageAccessToken: accessToken,
	}
}

func (h *RobotHandler) Handle(ctx context.Context, message *model.Message) error {
	if message.SentimentType == "" {
		message.SentimentType = db.MessageTemplateSentimentTypePositive // default value.
	}
	messageTemplate, err := h.Repo.QueryOne(ctx, message.SentimentType)
	if err != nil {
		return err
	}

	msgSender := new(messenger.Messenger)

	msg := messenger.NewTextMessage(
		message.GetRecipientID(),
		messageTemplate.Text,
	)
	messenger.HTTPClient = &MockFBClient{}
	return msgSender.SendMessage(ctx, msg, h.PageAccessToken)
}

type MockFBClient struct {
	http.Client
	Response []byte
	Err      error
}

func (c *MockFBClient) Do(req *http.Request) (*http.Response, error) {
	if c.Err != nil {
		return nil, c.Err
	}

	if len(c.Response) == 0 {
		c.Response = []byte("Bye~ ")
	}

	bodyReader, _ := req.GetBody()
	data, _ := io.ReadAll(bodyReader)

	log.Printf("MockFBClient has received: %s", string(data))

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(c.Response)),
	}, nil
}
