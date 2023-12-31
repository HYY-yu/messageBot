package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"messageBot/db/model"
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
	optFuns := make([]func(o *huggingface.InferenceClientOptions), 0)
	if !model.Prod {
		optFuns = append(optFuns, func(o *huggingface.InferenceClientOptions) {
			o.HTTPClient = &MockHuggingFaceClient{}
		})
	}
	ic := huggingface.NewInferenceClient(h.Token, optFuns...)
	log.Println("InferenceClient text is ", message.GetText())
	res, err := ic.ZeroShotClassification(ctx, &huggingface.ZeroShotClassificationRequest{
		Model:  "MoritzLaurer/DeBERTa-v3-base-mnli-fever-anli",
		Inputs: []string{message.GetText()},
		Parameters: huggingface.ZeroShotClassificationParameters{
			CandidateLabels: []string{"negative", "positive"},
		},
	})
	if err != nil {
		return err
	}

	if len(res) > 0 && len(res[0].Labels) == 2 {
		log.Printf("the huggingface res is %v \n", res)
		// the labels sorted in descending order.
		var negativeScore, positiveScore float64
		switch res[0].Labels[0] {
		case "negative":
			negativeScore = res[0].Scores[0]
			positiveScore = res[0].Scores[1]
		case "positive":
			positiveScore = res[0].Scores[0]
			negativeScore = res[0].Scores[1]
		}
		log.Printf("Negative score: %f, Positive score: %f \n", negativeScore, positiveScore)
		if negativeScore > positiveScore {
			message.SentimentType = model.MessageTemplateSentimentTypeNegative
		} else {
			message.SentimentType = model.MessageTemplateSentimentTypePositive
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
