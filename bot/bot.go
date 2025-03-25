package bot

import (
	"log"
	"strconv"
	"time"

	"github.com/aamarsel/browz/handlers"
	"gopkg.in/telebot.v3"
)

var Bot *telebot.Bot

func InitBot(token string) error {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		return err
	}

	Bot = b

	b.Handle(telebot.OnCallback, handlers.CallbackHandler)
	b.Handle("/start", handlers.StartHandler)
	b.Handle("/admin", func(c telebot.Context) error {
		log.Println("ID пользователя:", c.Sender().ID)
		return c.Send("Zuhr! Твой ID: " + strconv.FormatInt(c.Sender().ID, 10))
	})
	b.Handle("/book", handlers.BookHandler)
	b.Handle(telebot.OnText, handlers.MessageHandler)
	b.Handle(telebot.OnContact, handlers.ContactHandler)

	log.Println("🤖 Бот запущен!")
	b.Start()
	return nil
}
