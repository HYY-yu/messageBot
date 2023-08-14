package messenger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"messageBot/messenger/model"
	"net"
	"net/http"
	"time"
)

type MessageEntity interface {
	ToJson() []byte
}

type Messenger struct {
}

func (m *Messenger) SendMessage(ctx context.Context, message MessageEntity, pageAccessToken string) error {
	_, err := m.post(ctx, model.SendAPIEndpoint, message.ToJson(), pageAccessToken)
	if err != nil {
		return err
	}

	return nil
}

func (b *Messenger) post(ctx context.Context, url string, data []byte, token string) ([]byte, error) {
	url = fmt.Sprintf("%s?access_token=%s", url, token)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := HTTPClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		err := errors.New(fmt.Sprintf("Response code: %d. Body: %s", resp.StatusCode, string(body)))
		log.Println(err)
		return nil, err
	}

	return body, nil
}

type TextMessage struct {
	model.MessageResponse
	Text string `json:"text"`
}

func NewTextMessage(recipientID string, text string) *TextMessage {
	var t TextMessage
	t.MessageType = "RESPONSE"
	t.Recipient.ID = recipientID
	t.Text = text
	return &t
}

func (m *TextMessage) ToJson() []byte {
	b, _ := json.Marshal(m)
	return b
}

type HClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var HTTPClient HClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          50,
		MaxIdleConnsPerHost:   50,
		MaxConnsPerHost:       500,
		ForceAttemptHTTP2:     true,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}
