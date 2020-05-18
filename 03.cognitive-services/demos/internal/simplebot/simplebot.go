package simplebot

import (
	"log"

	"github.com/proge-software/uninapoli-csml-csbot/internal/tgconf"
	"github.com/proge-software/uninapoli-csml-csbot/pkg/tglog"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Bot WhySoSerious Bot implementation
type Bot struct {
	tbot tb.Bot
}

//NewFromEnv Bot constructor
func NewFromEnv() (*Bot, error) {
	c, err := tgconf.GetConfigurationsFromEnv()
	if err != nil {
		log.Fatalln(err)
	}

	bot, err := New(*c)
	if err != nil {
		return nil, err
	}

	return bot, nil
}

//New Bot constructor
func New(c tgconf.Configuration) (*Bot, error) {
	bot, err := new(c)
	if err != nil {
		return nil, err
	}

	{ // Handler: Text
		bot.tbot.Handle(tb.OnText, tglog.Wrap(bot.onText))
	}

	return bot, nil
}

func new(c tgconf.Configuration) (*Bot, error) {
	// telebot settings
	tbSettings := tb.Settings{
		Token: c.Token,
		Poller: &tb.LongPoller{
			Timeout: c.PollerInterval,
		},
	}

	// instantiate telebot bot
	tbot, err := tb.NewBot(tbSettings)
	if err != nil {
		return nil, err
	}

	// instantiate our bot
	bot := &Bot{
		tbot: *tbot,
	}
	return bot, nil
}

// Start starts the telegram bot
func (b *Bot) Start() {
	b.tbot.Start()
}
