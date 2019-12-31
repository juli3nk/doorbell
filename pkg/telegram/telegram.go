package telegram

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func Send(token string, userID int, text string) error {
	poller := &tb.LongPoller{Timeout: 15 * time.Second}

	settings := tb.Settings{
		Token:  token,
		Poller: poller,
	}

	bot, err := tb.NewBot(settings)
	if err != nil {
		return err
	}

	user := tb.User{
		ID: userID,
	}
	msg := &tb.Message{Sender: &user}

	bot.Send(msg.Sender, text)

	return  nil
}
