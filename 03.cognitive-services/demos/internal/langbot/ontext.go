package langbot

import (
	"context"
	"log"

	"github.com/proge-software/uninapoli-csml-csbot/pkg/wsssentiment"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *Bot) onText(m *tb.Message) (*tb.Message, error) {
	var message string
	ctx := context.Background()

	taChan := b.invokeTextAnalysisAPI(ctx, m.Text)

	taRes := <-taChan
	if taRes == nil || taRes.SentimentScore == nil {
		message = "Hello, unfortunately I'm not smart enough to understand your sentiments. Please provide me with empathy!"
	} else {
		score := taRes.SentimentScore
		if *score <= 0.3 { // Negative
			message += "I feel you a little bit upset!"
		} else if *score >= 0.7 { // Positive
			message += "What a nice positive message! Thanks!"
		} else { // Neutral
			message += "Why so serious?"
		}
	}

	return b.tbot.Send(m.Chat, message)
}

func (b *Bot) invokeTextAnalysisAPI(ctx context.Context, text string) chan *wsssentiment.TextAnalyticsResult {
	taChan := make(chan *wsssentiment.TextAnalyticsResult)
	if b.textAnalyticsCli == nil {
		defer close(taChan)
		return taChan
	}

	go func() {
		defer close(taChan)
		res, err := b.textAnalyticsCli.InvokeTextAnalytics(ctx, text)
		if err != nil {
			log.Printf("error invoking analytics service: %v", err)
		}
		taChan <- res
	}()

	return taChan
}
