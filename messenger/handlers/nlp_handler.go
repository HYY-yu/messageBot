package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"messageBot/db"
	"messageBot/messenger/model"
	"net/http"

	huggingface "github.com/hupe1980/go-huggingface"
)

type NLPHandler struct {
	Token string
}

func NewNLPHandler(token string) *NLPHandler {
	return &NLPHandler{
		Token: token,
	}
}

func (h *NLPHandler) Handle(ctx context.Context, message *model.Message) error {
	ic := huggingface.NewInferenceClient(h.Token, func(o *huggingface.InferenceClientOptions) {
		o.HTTPClient = &MockHuggingFaceClient{}
	})
	res, err := ic.ZeroShotClassification(ctx, &huggingface.ZeroShotClassificationRequest{
		Model:  "facebook/bart-large-mnli",
		Inputs: []string{message.GetText()},
		Parameters: huggingface.ZeroShotClassificationParameters{
			CandidateLabels: []string{"positive", "negative"},
		},
	})
	if err != nil {
		return err
	}

	if len(res) > 0 && len(res[0].Labels) == 2 {
		// the labels sorted in descending order.
		negativeScore := res[0].Scores[0]
		positiveScore := res[0].Scores[1]
		if negativeScore > positiveScore {
			message.SentimentType = db.MessageTemplateSentimentTypeNegative
		} else {
			message.SentimentType = db.MessageTemplateSentimentTypePositive
		}
	}
	return nil
}

type MockHuggingFaceClient struct {
	Err error
}

func (c *MockHuggingFaceClient) Do(req *http.Request) (*http.Response, error) {
	if c.Err != nil {
		return nil, c.Err
	}
	data := []float64{0.9, 0.1}
	if rand.Intn(2) == 1 {
		data = []float64{0.1, 0.9}
	}

	resp := []struct {
		Sequence string    `json:"sequence,omitempty"`
		Labels   []string  `json:"labels,omitempty"`
		Scores   []float64 `json:"scores,omitempty"`
	}{
		{
			Sequence: "TEST",
			Labels:   []string{"negative", "positive"},
			Scores:   data,
		},
	}
	resultText := "negative"
	if data[0] < data[1] {
		resultText = "positive"
	}
	log.Println("response data is: ", resultText)

	respBytes, _ := json.Marshal(resp)

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(respBytes)),
	}, nil
}
