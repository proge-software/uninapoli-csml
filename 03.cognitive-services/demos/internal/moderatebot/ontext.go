package moderatebot

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/proge-software/uninapoli-csml-csbot/pkg/wssmoderator"
	"github.com/proge-software/uninapoli-csml-csbot/pkg/wsssentiment"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *Bot) onText(m *tb.Message) (*tb.Message, error) {
	var message string
	ctx := context.Background()

	taChan := b.invokeTextAnalysisAPI(ctx, m.Text)
	mChan := b.invokeModeratorTextAPI(ctx, m.Text)

	taRes := <-taChan
	if taRes == nil || taRes.SentimentScore == nil {
		message = "Hello, unfortunately I'm not smart enough to understand your sentiments. Please provide me with empathy!"
	} else {
		score := taRes.SentimentScore
		if *score <= 0.3 {
			message += "I feel you a little bit upset!"
		} else if *score >= 0.7 {
			message += "What a nice positive message! Thanks!"
		} else {
			message += "Why so serious?"
		}
	}

	mRes := <-mChan
	if mRes != nil && len(mRes.BadWords) > 0 {
		badlist := strings.Join(mRes.BadWords, ", ")
		message += fmt.Sprintf("\nHowever, you should avoid to use bad words like: %s.", badlist)
	}

	return b.tbot.Send(m.Chat, message)
}

func (b *Bot) invokeModeratorTextAPI(ctx context.Context, text string) chan *wssmoderator.ContentModeratorTextResult {
	mChan := make(chan *wssmoderator.ContentModeratorTextResult, 1)
	if b.moderatorCli == nil {
		close(mChan)
		return mChan
	}

	go func() {
		defer close(mChan)
		moderatorResult, err := b.moderatorCli.InvokeContentModeratorText(ctx, text)
		if err != nil {
			log.Printf("error invoking moderator service: %v", err)
		}
		mChan <- moderatorResult
	}()

	return mChan
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
