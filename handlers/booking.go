package handlers

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aamarsel/browz/database"
	"github.com/aamarsel/browz/keyboards"
	"github.com/aamarsel/browz/utils"
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
	parsedTime, _ := time.Parse("15:04", parts[2][11:16])
	formattedTime := parsedTime.Format("15:04:00")
	userID := c.Sender().ID

	slotID, err := database.FindAvailableSlot(formattedDate, formattedTime)
	if err != nil {
		return c.Send("Слот уже забронирован. Выберите другой.")
	}

	err = database.BookSlot(userID, slotID, serviceID)
	if err != nil {
		log.Println("Ошибка при бронировании слота:", err)
		return c.Send("Ошибка при создании бронирования.")
	}

	// Получаем данные о клиенте
	client, err := database.GetClientByID(userID)
	if err != nil {
		log.Println("Ошибка при получении данных клиента:", err)
		return c.Send("Произошла ошибка, попробуйте позже.")
	}

	// Получаем данные об услуге
	service, err := database.GetServiceByID(serviceID)
	if err != nil {
		log.Println("Ошибка при получении данных об услуге:", err)
		return c.Send("Произошла ошибка, попробуйте позже.")
	}

	// Telegram ID Зухры
	zukhraID, err := strconv.ParseInt(os.Getenv("ZUKHRA_ID"), 10, 64)
	if err != nil {
		log.Println("Ошибка при парсинге ZUKHRA_ID:", err)
		return c.Send("Ошибка конфигурации.")
	}

	// Формируем красивое сообщение для Зухры
	zukhraMessage := fmt.Sprintf(
		"📅 *Новая запись!*\n\n"+
			"👤 *Клиент:* %s\n"+
			"📞 *Телефон:* %s\n"+
			"🛠 *Услуга:* %s\n"+
			"📆 *Дата:* %s\n"+
			"⏰ *Время:* %s\n\n"+
			"🔔 Подтвердите или отклоните запись.",
		client.Name,
		client.Phone,
		service.Name,
		date,
		parts[2][11:16], // Отображаем время в исходном формате
	)

	// Отправляем сообщение Зухре
	zukhraUser := &telebot.User{ID: zukhraID}
	c.Bot().Send(zukhraUser, zukhraMessage, telebot.ModeMarkdown)

	// Отправляем подтверждение пользователю
	return keyboards.SendMainMenu(c, "✅ Ваша запись отправлена на подтверждение мастеру Зухре. Ожидайте, вам придет уведомление.")
}

// HandleMyBookings обрабатывает нажатие на кнопку "📅 Мои бронирования"
func HandleMyBookings(c telebot.Context) error {
	clientID := c.Sender().Recipient() // Telegram ID пользователя

	// Получаем все бронирования пользователя
	bookings, err := database.GetUserBookings(clientID)
	if err != nil {
		log.Println("Ошибка при получении бронирований:", err)
		return c.Send("Ошибка при загрузке бронирований. Попробуйте позже.")
	}

	// Если записей нет
	if len(bookings) == 0 {
		return c.Send("У вас пока нет бронирований.")
	}

	// Отправляем каждую запись отдельным сообщением
	for _, booking := range bookings {
		ratingText := "Нет оценки"
		if booking.Rating != nil {
			ratingText = fmt.Sprintf("%d/5", *booking.Rating)
		}
		// Форматируем сообщение с записями
		msgText := fmt.Sprintf(
			"📅 *Дата:* %s\n"+
				"💆 *Услуга:* %s\n"+
				"💵 *Цена:* %d ₽\n"+
				"⭐ *Оценка:* %s\n"+
				"🔹 *Статус:* %s",
			booking.DateTime.Format("02.01.2006 15:04"),
			booking.ServiceName,
			booking.ServicePrice,
			ratingText,
			utils.FormatStatus(booking.Status),
		)

		// Проверяем, можно ли отменить запись (если она в будущем)
		btns := &telebot.ReplyMarkup{}
		if booking.DateTime.After(time.Now()) && (booking.Status != "cancelled" && booking.Status != "completed") {
			cancelBtn := btns.Data("❌ Отменить запись", "cancel_booking", booking.ID)
			btns.Inline(btns.Row(cancelBtn))
		}

		// Отправляем сообщение
		c.Send(msgText, btns, telebot.ModeMarkdown)
	}

	return nil
}

// HandleCancelBooking обрабатывает кнопку "❌ Отменить запись"
func HandleCancelBooking(c telebot.Context) error {
	telegramID := c.Sender().Recipient() // Получаем Telegram ID пользователя
	// Получаем bookingID из callback данных
	bookingID := strings.Split(c.Data(), "|")[1]

	// Отменяем бронирование
	err := database.CancelBooking(telegramID, bookingID)
	if err != nil {
		log.Println("Ошибка при отмене бронирования:", err)
		return c.Respond(&telebot.CallbackResponse{
			Text:      "Не удалось отменить запись. Возможно, она уже отменена или завершена.",
			ShowAlert: true,
		})
	}

	// Отвечаем пользователю
	return c.Respond(&telebot.CallbackResponse{
		Text:      "✅ Бронирование успешно отменено!",
		ShowAlert: true,
	})
}
