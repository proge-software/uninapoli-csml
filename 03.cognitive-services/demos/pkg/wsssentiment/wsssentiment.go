package wsssentiment

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.0/textanalytics"
	"github.com/Azure/go-autorest/autorest"
)

//TextAnalyticsServiceClient client for the Azure TextAnalytics Service
type TextAnalyticsServiceClient struct {
	conf             Configuration
	textanalyticsCli *textanalytics.BaseClient
}

//NewTextAnalyticsServiceClient TextAnalyticsServiceClient constructor
func NewTextAnalyticsServiceClient(conf *Configuration) *TextAnalyticsServiceClient {
	if !conf.IsValid() {
		return nil
	}

	client := textanalytics.New(conf.ServiceEnpoint)
	client.Authorizer = autorest.NewCognitiveServicesAuthorizer(conf.TextAnalyticsSubscription)

	return &TextAnalyticsServiceClient{
		conf:             *conf,
		textanalyticsCli: &client,
	}
}

//InvokeTextAnalytics invokes the TextAnalytics APIs with the provided photo
func (s *TextAnalyticsServiceClient) InvokeTextAnalytics(textAnalyticsContext context.Context, message string) (*TextAnalyticsResult, error) {
	if s == nil {
		return nil, fmt.Errorf("text analytics service client is not initialized")
	}

	id, language := "singledoc", "en"
	messages := []textanalytics.MultiLanguageInput{{ID: &id, Text: &message, Language: &language}}
	res, err := s.textanalyticsCli.Sentiment(textAnalyticsContext,
		textanalytics.MultiLanguageBatchInput{Documents: &messages})
	if err != nil {
		return nil, err
	}

	docs := res.Documents
	if docs == nil || len(*docs) == 0 {
		return &TextAnalyticsResult{SentimentScore: nil}, nil
	}
	return &TextAnalyticsResult{SentimentScore: (*docs)[0].Score}, nil
}

// TextAnalyticsResult result of the TextAnalyticsAPI
type TextAnalyticsResult struct {
	SentimentScore *float64
}
