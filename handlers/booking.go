package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/aamarsel/browz/database"
	"github.com/aamarsel/browz/keyboards"
	"gopkg.in/telebot.v3"
)

func BookHandler(c telebot.Context) error {
	return keyboards.ShowDatePicker(c)
}

func ConfirmBookingHandler(c telebot.Context) error {
	parts := strings.Split(c.Data(), "|")
	if len(parts) < 3 {
		return c.Send("Ошибка при обработке запроса.")
	}

	serviceID, _ := strconv.Atoi(parts[1])
	date := parts[2][:10]
	parsedDate, _ := time.Parse("02.01.2006", date)
	formattedDate := parsedDate.Format("2006-01-02")
	parsedTime, _ := time.Parse("15:04", parts[2][12:17])
	formattedTime := parsedTime.Format("15:04:00")
	userID := c.Sender().ID

	slotID, err := database.FindAvailableSlot(formattedDate, formattedTime)
	if err != nil {
		return c.Send("Слот уже забронирован. Выберите другой.")
	}

	err = database.BookSlot(userID, slotID, serviceID)
	if err != nil {
		return c.Send("Ошибка при создании бронирования.")
	}

	return c.Send("✅ Ваша запись отправлена на подтверждение мастеру Зухре. Ожидайте, вам придет уведомление.", keyboards.MainMenu)
}
