package handlers

import (
	"github.com/aamarsel/browz/keyboards"
	"github.com/aamarsel/browz/models"
	"gopkg.in/telebot.v3"
)

func MessageHandler(c telebot.Context) error {
	userID := c.Sender().ID

	switch models.UserState[userID] {
	case models.StateAwaitingName:
		return ProcessNameInput(c)
	case models.StateAwaitingContact:
		return c.Send("Пожалуйста, отправьте контакт кнопкой ниже.")
	}

	switch c.Text() {
	case "➕ Записаться к Зухре":
		return BookHandler(c)
	default:
		return c.Send("Я не понял команду.", keyboards.MainMenu)
	}
}
