package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aamarsel/browz/database"
	"github.com/aamarsel/browz/keyboards"
	"github.com/aamarsel/browz/models"
	"gopkg.in/telebot.v3"
)

func HandleEditServices(c telebot.Context) error {
	services, err := database.GetAllServices()
	if err != nil {
		log.Println("Ошибка получения услуг:", err)
		return c.Send("Ошибка получения списка услуг.")
	}

	if len(services) == 0 {
		return c.Send("Нет добавленных услуг.")
	}

	for _, service := range services {
		markup := &telebot.ReplyMarkup{}
		btnEdit := markup.Data("✏️ Редактировать", "edit_service", strconv.Itoa(service.ID))
		btnDelete := markup.Data("❌ Удалить", "delete_service", strconv.Itoa(service.ID))
		markup.Inline(markup.Row(btnEdit, btnDelete))

		durationMinutes := int(service.Duration.Minutes())

		text := fmt.Sprintf(
			"🛠 *%s*\n💰 Цена: %d ₽\n⏳ Длительность: %d мин",
			service.Name,
			service.Price,
			durationMinutes,
		)

		c.Send(text, markup, telebot.ModeMarkdown)
	}

	return nil
}

func HandleDeleteService(c telebot.Context) error {
	serviceID, _ := strconv.Atoi(strings.Split(c.Data(), "|")[1])

	err := database.DeleteService(serviceID)
	if err != nil {
		log.Println("Ошибка при удалении услуги:", err)
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка удаления."})
	}

	return c.Respond(&telebot.CallbackResponse{Text: "✅ Услуга удалена."})
}

func HandleEditService(c telebot.Context) error {
	serviceID, _ := strconv.Atoi(strings.Split(c.Data(), "|")[1])
	userID := c.Sender().ID
	models.TempServiceData[userID] = models.TempService{ID: serviceID}
	models.UserState[userID] = models.StateAwaitingServiceNameEdit

	c.Send("Введите новое название услуги:")

	return nil
}

func HandleAddService(c telebot.Context) error {
	userID := c.Sender().ID
	models.UserState[userID] = models.StateAwaitingServiceName

	return c.Send("Введите название услуги:")
}

func ProcessServiceName(c telebot.Context) error {
	userID := c.Sender().ID
	name := c.Text()

	models.TempServiceData[userID] = models.TempService{Name: name}
	models.UserState[userID] = models.StateAwaitingServicePrice

	return c.Send("Введите цену услуги (в рублях):")
}

func ProcessServiceNameEdit(c telebot.Context) error {
	userID := c.Sender().ID
	name := c.Text()

	tempService := models.TempServiceData[userID]
	tempService.Name = name
	models.TempServiceData[userID] = tempService
	models.UserState[userID] = models.StateAwaitingServicePriceEdit

	return c.Send("Введите новую цену услуги (в рублях):")
}

func ProcessServicePriceEdit(c telebot.Context) error {
	userID := c.Sender().ID
	price, err := strconv.Atoi(c.Text())
	if err != nil || price < 0 {
		return c.Send("Введите корректную цену (число в рублях):")
	}

	tempService := models.TempServiceData[userID]
	tempService.Price = price
	models.TempServiceData[userID] = tempService
	models.UserState[userID] = models.StateAwaitingServiceDurationEdit

	return c.Send("Введите новую длительность услуги в минутах:")
}

func ProcessServiceDurationEdit(c telebot.Context) error {
	userID := c.Sender().ID
	duration, err := strconv.Atoi(c.Text())
	if err != nil || duration <= 0 {
		return c.Send("Введите корректную длительность (число в минутах):")
	}

	tempService := models.TempServiceData[userID]
	tempService.Duration = duration

	delete(models.TempServiceData, userID)
	delete(models.UserState, userID)

	err = database.EditService(tempService)
	if err != nil {
		log.Println("Ошибка при редактировании услуги:", err)
		return c.Send("Ошибка при редактировании услуги.")
	}

	return c.Send("✅ Услуга успешно отредактирована.")
}

func ProcessServicePrice(c telebot.Context) error {
	userID := c.Sender().ID
	price, err := strconv.Atoi(c.Text())
	if err != nil || price < 0 {
		return c.Send("Введите корректную цену (число в рублях):")
	}

	models.TempServiceData[userID] = models.TempService{
		Name:  models.TempServiceData[userID].Name,
		Price: price,
	}
	models.UserState[userID] = models.StateAwaitingServiceDuration

	return c.Send("Введите длительность услуги в минутах:")
}

func ProcessServiceDuration(c telebot.Context) error {
	userID := c.Sender().ID
	duration, err := strconv.Atoi(c.Text())
	if err != nil || duration <= 0 {
		return c.Send("Введите корректную длительность (число в минутах):")
	}

	tempService := models.TempServiceData[userID]
	tempService.Duration = duration

	err = database.AddService(tempService)
	if err != nil {
		log.Println("Ошибка при добавлении услуги:", err)
		return keyboards.SendMainMenu(c, "Ошибка при добавлении услуги.")
	}

	delete(models.TempServiceData, userID)
	delete(models.UserState, userID)

	return keyboards.SendMainMenu(c, "✅ Услуга успешно добавлена!")
}
