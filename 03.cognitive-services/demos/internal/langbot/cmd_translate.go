package langbot

import (
	"context"
	"log"
	"regexp"

	"github.com/proge-software/uninapoli-csml-csbot/pkg/wsstranslator"
	tb "gopkg.in/tucnak/telebot.v2"
)

const langPattern = "-to:(.*?) "

func (b *Bot) translate(m *tb.Message) (*tb.Message, error) {
	if !m.Private() {
		return nil, nil
	}

	language := extractLang(m.Payload)
	message := m.Payload
	if language != nil {
		slack := 4 + len(*language) + 1
		message = message[slack:]
	}

	ctx := context.Background()
	res, err := b.translatorCli.InvokeTranslator(ctx, message, language)
	if err != nil {
		log.Printf("error invoking translator service: %v", err)
		return nil, err
	}

	response := generateTranslationResponse(res)
	return b.tbot.Send(m.Chat, response)
}

func generateTranslationResponse(res *wsstranslator.TranslatorResult) string {
	var message string
	if res.IdentifiedLang != nil {
		message += "Identified language " + *res.IdentifiedLang + ". "
	}

	if res.Translation != nil {
		message += "Translation:\n\n" + *res.Translation
	} else {
		message += "Unfortunately, I was not able to translate what you said"
	}

	return message
}

func extractLang(payload string) *string {
	exp, err := regexp.Compile(langPattern)
	if err != nil {
		log.Printf("can not compile regexp: %s", langPattern)
		return nil
	}

	if submatches := exp.FindStringSubmatch(payload); len(submatches) >= 1 {
		return &submatches[1]
	}
	return nil
}
