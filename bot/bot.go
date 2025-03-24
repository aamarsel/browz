package bot

import (
	"log"
	"strconv"
	"time"

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

	InitKeyboards()

	b.Handle(telebot.OnCallback, CallbackHandler)
	b.Handle("/start", StartHandler)
	b.Handle("/admin", func(c telebot.Context) error {
		log.Println("ID пользователя:", c.Sender().ID)
		return c.Send("Zuhr! Твой ID: " + strconv.FormatInt(c.Sender().ID, 10))
	})
	b.Handle("/book", BookHandler)
	b.Handle("/appointments", ListAppointments)
	b.Handle(telebot.OnText, MessageHandler)
	b.Handle(telebot.OnContact, ContactHandler)

	log.Println("🤖 Бот запущен!")
	b.Start()
	return nil
}
