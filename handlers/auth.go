package handlers

import (
	"log"

	"github.com/aamarsel/browz/database"
	"github.com/aamarsel/browz/keyboards"
	"github.com/aamarsel/browz/models"
	"gopkg.in/telebot.v3"
)

func StartHandler(c telebot.Context) error {
	userID := c.Sender().ID

	// Проверяем, есть ли клиент в БД
	exists, err := database.ClientExists(userID)
	if err != nil {
		log.Println("Ошибка проверки клиента:", err)
		return c.Send("Произошла ошибка. Попробуйте позже.")
	}

	if exists {
		keyboards.SendMainMenu(c, "👋 С возвращением!")
		return nil
	}

	// Начинаем регистрацию
	models.RegistrationStorage[userID] = models.RegistrationState{}
	models.UserState[userID] = models.StateAwaitingName

	return c.Send("👋 Привет! Введите ваше имя:")
}

func ProcessNameInput(c telebot.Context) error {
	userID := c.Sender().ID
	name := c.Text()

	state := models.RegistrationStorage[userID]
	state.Name = name
	models.RegistrationStorage[userID] = state
	models.UserState[userID] = models.StateAwaitingContact

	contactBtn := telebot.ReplyButton{Text: "📱 Поделиться контактом", Contact: true}
	menu := &telebot.ReplyMarkup{
		ReplyKeyboard:   [][]telebot.ReplyButton{{contactBtn}},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	return c.Send("Теперь отправьте свой контакт:", menu)
}

func ContactHandler(c telebot.Context) error {
	userID := c.Sender().ID

	if c.Message().Contact == nil {
		return c.Send("Пожалуйста, используйте кнопку ниже для отправки номера.")
	}

	phone := c.Message().Contact.PhoneNumber
	state := models.RegistrationStorage[userID]
	state.Phone = phone

	err := database.SaveClient(userID, state.Name, state.Phone)
	if err != nil {
		log.Println("Ошибка сохранения клиента:", err)
		return c.Send("Ошибка при регистрации. Попробуйте позже.")
	}

	delete(models.RegistrationStorage, userID)
	models.UserState[userID] = models.StateNone

	keyboards.SendMainMenu(c, "✅ Регистрация завершена!")
	return nil
}
