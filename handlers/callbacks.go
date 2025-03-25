package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aamarsel/browz/database"
	"github.com/aamarsel/browz/keyboards"
	"github.com/aamarsel/browz/models"
	"github.com/aamarsel/browz/utils"
	"gopkg.in/telebot.v3"
)

func CallbackHandler(c telebot.Context) error {
	callbackData := c.Callback().Data

	if strings.Contains(callbackData, "pick_date") {
		return DatePickerHandler(c)
	} else if strings.Contains(callbackData, "pick_slot") {
		return SlotPickerHandler(c)
	} else if strings.Contains(callbackData, "pick_service") {
		return ServicePickerHandler(c)
	} else if strings.Contains(callbackData, "confirm_booking") {
		return ConfirmBookingHandler(c)
	} else if strings.Contains(callbackData, "cancel_booking") {
		return HandleCancelBooking(c)
	} else {
		log.Println("Ошибка! Неизвестный callback:", callbackData)
	}

	return nil
}

func DatePickerHandler(c telebot.Context) error {
	parts := strings.Split(c.Data(), "|")
	if len(parts) < 2 {
		return c.Send("Ошибка: неверный формат данных.")
	}

	selectedDate := parts[1]
	return keyboards.ShowTimeSlots(c, selectedDate)
}

func SlotPickerHandler(c telebot.Context) error {
	data := c.Data()
	userID := c.Sender().ID

	models.TempStorage[userID] = models.SelectedSlot{
		Date: data[:10],
		Time: data[11:],
	}

	// Показываем список услуг
	return c.Send("Выберите услугу:", &telebot.ReplyMarkup{
		InlineKeyboard: keyboards.GetServicesButtons(),
	})
}

func ServicePickerHandler(c telebot.Context) error {
	data := c.Data()
	serviceID, err := strconv.Atoi(strings.TrimPrefix(data, "pick_service:"))
	if err != nil {
		log.Println(err)
		return c.Send("Ошибка при выборе услуги. Попробуйте снова.")
	}

	var name string
	var price int
	var duration time.Duration
	err = database.DB.QueryRow(
		context.Background(),
		"SELECT name, price, duration FROM services WHERE id = $1",
		serviceID,
	).Scan(&name, &price, &duration)
	if err != nil {
		return c.Send("Ошибка: услуга не найдена.")
	}

	userID := c.Sender().ID
	slot, exists := models.TempStorage[userID]
	if !exists {
		return c.Send("Ошибка: выберите слот перед выбором услуги.")
	}

	timeStr := slot.Time[11:]

	msg := fmt.Sprintf(
		"📅 Дата: %s\n💆‍♀️ Услуга: %s\n⏳ Время работы мастера: %s\n💰 Цена: %d руб\n\n"+
			"Подтвердите свою запись к своему любимому мастеру Зухре 😊",
		timeStr,
		name,
		utils.FormatDuration(duration),
		price,
	)

	btnYes := telebot.InlineButton{Text: "✅ Да", Data: fmt.Sprintf("confirm_booking|%d|%s", serviceID, timeStr)}
	btnNo := telebot.InlineButton{Text: "❌ Нет", Data: "cancel_booking"}

	return c.Send(msg, &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{{btnYes}, {btnNo}},
	})
}
