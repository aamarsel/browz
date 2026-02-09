package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aamarsel/browz/database"
	"gopkg.in/telebot.v3"
)

func HandleFutureBookings(c telebot.Context) error {
	appointments, err := database.GetBookingsByStatus("accepted", true)
	if err != nil {
		return c.Send("Ошибка при получении записей.")
	}

	if len(appointments) == 0 {
		return c.Send("Нет будущих записей.")
	}

	for _, booking := range appointments {
		msgText := fmt.Sprintf(
			"📅 *Дата:* %s\n"+
				"💆 *Услуга:* %s\n"+
				"👤 *Клиент:* %s\n",
			booking.DateTime.Format("02.01.2006 15:04"),
			booking.ServiceName,
			booking.ClientName,
		)

		btns := &telebot.ReplyMarkup{}
		cancelBtn := btns.Data("❌ Отменить запись", "cancel_admin_booking", booking.ID)
		btns.Inline(btns.Row(cancelBtn))

		c.Send(msgText, btns, telebot.ModeMarkdown)
	}
	return nil
}

func HandleOldBookings(c telebot.Context) error {
	bookings, err := database.GetPastBookings(c)
	if err != nil {
		return c.Send("Ошибка при получении записей.")
	}

	for _, booking := range bookings {
		msgText := fmt.Sprintf(
			"📅 *Дата:* %s\n"+
				"💆 *Услуга:* %s\n"+
				"👤 *Клиент:* %s\n",
			booking.DateTime.Format("02.01.2006 15:04"),
			booking.ServiceName,
			booking.ClientName,
		)

		c.Send(msgText, telebot.ModeMarkdown)
	}
	return nil
}

func HandlePendingBookings(c telebot.Context) error {
	appointments, err := database.GetBookingsByStatus("pending", false)
	if err != nil {
		return c.Send("Ошибка при получении записей.")
	}

	if len(appointments) == 0 {
		return c.Send("Нет неподтвержденных записей.")
	}

	for _, booking := range appointments {
		msgText := fmt.Sprintf(
			"📅 *Дата:* %s\n"+
				"💆 *Услуга:* %s\n"+
				"👤 *Клиент:* %s\n",
			booking.DateTime.Format("02.01.2006 15:04"),
			booking.ServiceName,
			booking.ClientName,
		)

		btns := &telebot.ReplyMarkup{}
		acceptBtn := btns.Data("✅ Принять", "accept_booking", booking.ID)
		declineBtn := btns.Data("❌ Отклонить", "decline_booking", booking.ID)

		btns.Inline(btns.Row(acceptBtn, declineBtn))

		c.Send(msgText, btns, telebot.ModeMarkdown)
	}
	return nil
}

func HandleAcceptBooking(c telebot.Context) error {
	bookingID := strings.Split(c.Data(), "|")[1]

	// Обновляем статус бронирования
	booking, err := database.UpdateBookingStatus(bookingID, "accepted")
	if err != nil {
		return c.Send("Ошибка при подтверждении записи.")
	}

	// Формируем красивое уведомление для клиента
	notification := fmt.Sprintf(
		"✅ *Ваша запись подтверждена!*\n\n"+
			"*Дата:* %s\n"+
			"*Время:* %s\n"+
			"*Услуга:* %s\n\n"+
			"📍 Ждем вас вовремя!",
		booking.DateTime.Format("02.01.2006"),
		booking.DateTime.Format("15:04"),
		booking.ServiceName,
	)

	// Отправляем уведомление пользователю
	id, _ := strconv.ParseInt(booking.ClientTelegramID, 10, 64)
	recipient := &telebot.User{ID: id}
	_, err = c.Bot().Send(recipient, notification, telebot.ModeMarkdown)
	if err != nil {
		log.Println("Ошибка при отправке уведомления пользователю:", err)
	}

	// Подтверждаем администратору
	return c.Send("Запись подтверждена ✅")
}

func HandleDeclineBooking(c telebot.Context) error {
	bookingID := strings.Split(c.Data(), "|")[1]
	_, err := database.UpdateBookingStatus(bookingID, "cancelled")
	if err != nil {
		return c.Send("Ошибка при отклонении записи.")
	}
	return c.Send("Запись отклонена ❌")
}
