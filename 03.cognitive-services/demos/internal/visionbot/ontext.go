package visionbot

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *Bot) onText(m *tb.Message) (*tb.Message, error) {
	message := "Hello man!"
	return b.tbot.Send(m.Chat, message)
}
