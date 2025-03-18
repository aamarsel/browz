package main

import (
	"log"
	"os"

	"github.com/aamarsel/browz/bot"
	"github.com/aamarsel/browz/database"
)

func main() {
	LoadConfig()

	err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()

	// Запускаем бота
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("❌ TELEGRAM_TOKEN не найден!")
	}

	err = bot.InitBot(token)

	if err != nil {
		log.Fatal(err)
	}
}
