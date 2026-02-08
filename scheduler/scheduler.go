package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
	"gopkg.in/telebot.v3"
)

// Функция для завершения записей

func StartScheduler(bot *telebot.Bot) {
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(15).Minutes().StartAt(time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)).Do(func() {
		CompleteOldBookings(bot)
	})
	scheduler.StartAsync()
}
