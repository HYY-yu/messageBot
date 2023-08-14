package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_fakeMessageTemplate_QueryOne(t *testing.T) {
	type args struct {
		ctx           context.Context
		sentimentType MessageTemplateSentimentType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TESTNegative",
			args: args{
				ctx:           context.Background(),
				sentimentType: MessageTemplateSentimentTypeNegative,
			},
			wantErr: false,
		},
		{
			name: "TESTPositive",
			args: args{
				ctx:           context.Background(),
				sentimentType: MessageTemplateSentimentTypePositive,
			},
			wantErr: false,
		},
		{
			name: "TESTErr",
			args: args{
				ctx:           context.Background(),
				sentimentType: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &fakeMessageTemplate{}
			got, err := o.QueryOne(tt.args.ctx, tt.args.sentimentType)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.args.sentimentType == MessageTemplateSentimentTypeNegative {
				assert.Contains(t, NegativeMessages, got.Text)
			} else if tt.args.sentimentType == MessageTemplateSentimentTypePositive {
				assert.Contains(t, PositiveMessages, got.Text)
			}
		})
	}
}
