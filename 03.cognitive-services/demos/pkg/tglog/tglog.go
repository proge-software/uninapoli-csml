package tglog

import (
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

//Wrap ....
func Wrap(onText func(*tb.Message) (*tb.Message, error)) func(*tb.Message) {
	return func(m *tb.Message) {

		if m == nil {
			log.Printf("[%d] empty message received", m.ID)
		} else {
			log.Printf("[%d] message received: %v", m.ID, m.Text)
		}

		resp, err := onText(m)
		if err != nil {
			log.Printf("[%d] error sending response: %v", m.ID, err)
		}

		if resp == nil {
			log.Printf("[%d] no response sent", m.ID)
		} else {
			log.Printf("[%d] response sent: %v", m.ID, resp.Text)
		}
	}
}
