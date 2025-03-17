package main

import (
	"log"
	"os"

	"time"

	"gopkg.in/telebot.v3"
)

func main() {
	// Загружаем токен бота из переменных окружения
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN не найден в переменных окружения")
	}

	// Настройки бота
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	// Создаем бота
	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	// Обрабатываем команду /start
	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send("Привет! Я бот для записи на брови. Выберите действие.")
	})

	log.Println("Бот запущен!")
	bot.Start()
}
