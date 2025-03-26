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
	case models.StateAwaitingServiceName:
		return ProcessServiceName(c)
	case models.StateAwaitingServiceDuration:
		return ProcessServiceDuration(c)
	case models.StateAwaitingServicePrice:
		return ProcessServicePrice(c)
	}

	switch c.Text() {
	case "➕ Записаться к Зухре":
		return BookHandler(c)
	case "📅 Мои бронирования":
		return HandleMyBookings(c)
	case "📅 Будущие записи":
		return HandleFutureBookings(c)
	case "⏳ Неподтвержденные записи":
		return HandlePendingBookings(c)
	case "✏️ Редактировать услуги":
		return HandleEditServices(c)
	case "➕ Добавить услугу":
		return HandleAddService(c)
	default:
		return keyboards.SendMainMenu(c, "Я не понял команду.")
	}
}
