package langbot

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/proge-software/uninapoli-csml-csbot/pkg/wsstranslator"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *Bot) languages(m *tb.Message) (*tb.Message, error) {
	if !m.Private() {
		return nil, nil
	}

	ctx := context.Background()
	res, err := b.translatorCli.InvokeLanguages(ctx)
	if err != nil {
		log.Printf("error invoking translator service: %v", err)
		return nil, err
	}

	pRespStr := b.langsMarkdown(res.Languages)
	return b.tbot.Send(m.Chat, *pRespStr, tb.ModeMarkdown)
}

func (b *Bot) langsMarkdown(languages []wsstranslator.Language) *string {
	var sb strings.Builder
	sort.Slice(languages, func(i, j int) bool {
		return languages[i].Code < languages[j].Code
	})

	for _, l := range languages {
		langline := fmt.Sprintf("*%s*: %s\n", l.Code, l.Name)
		sb.WriteString(langline)
	}
	str := sb.String()
	return &str
}
