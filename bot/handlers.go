package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/aamarsel/browz/auth"
	"github.com/aamarsel/browz/database"
)

func StartHandler(c telebot.Context) error {
	userID := c.Sender().ID

	var exists bool
	err := database.DB.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM clients WHERE telegram_id = $1)", userID).Scan(&exists)

	if err != nil {
		log.Println("Ошибка проверки клиента:", err)
		return c.Send("Произошла ошибка. Попробуйте позже.")
	}

	if exists {
		// Если клиент уже зарегистрирован, открываем главное меню
		return c.Send("👋 С возвращением!", MainMenu)
	}

	// Если клиента нет, начинаем регистрацию
	registrationStorage[userID] = RegistrationState{Name: "", Phone: ""}
	userState[userID] = StateAwaitingName
	return c.Send("👋 Привет!\nЭто бот для записи на брови к мастеру Зухре.\nДавай зарегистрируемся.\nКак тебя зовут?")
}

func processNameInput(c telebot.Context) error {
	userID := c.Sender().ID
	name := c.Text()

	// Сохраняем имя во временное хранилище
	state := registrationStorage[userID]
	state.Name = name
	registrationStorage[userID] = state

	// Меняем состояние пользователя
	userState[userID] = StateAwaitingContact

	// Отправляем кнопку для отправки контакта
	contactBtn := telebot.ReplyButton{Text: "📱 Поделиться контактом", Contact: true}
	menu := &telebot.ReplyMarkup{
		ReplyKeyboard:   [][]telebot.ReplyButton{{contactBtn}},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
	return c.Send("Теперь отправьте свой контакт, нажав на кнопку ниже:", menu)
}

func ContactHandler(c telebot.Context) error {
	userID := c.Sender().ID

	if c.Message().Contact == nil {
		return c.Send("Пожалуйста, отправьте свой номер через кнопку ниже.")
	}

	phone := c.Message().Contact.PhoneNumber
	state := registrationStorage[userID]
	state.Phone = phone

	// Сохраняем клиента в базе
	_, err := database.DB.Exec(context.Background(),
		"INSERT INTO clients (telegram_id, name, phone) VALUES ($1, $2, $3)",
		userID, state.Name, state.Phone)

	if err != nil {
		log.Println("Ошибка сохранения клиента:", err)
		return c.Send("Произошла ошибка при регистрации. Попробуйте позже.")
	}

	// Удаляем временные данные
	delete(tempStorage, userID)

	return c.Send("✅ Регистрация завершена! Добро пожаловать!", MainMenu)
}

func BookHandler(c telebot.Context) error {
	return showDatePicker(c)
}

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
		return BookHandler(c)
	} else {
		log.Println("Ошибка! Неизвестный callback:", callbackData)
	}

	return nil
}

func ListAppointments(c telebot.Context) error {
	if !auth.IsAdmin(c.Sender().ID) {
		return c.Send("❌ У вас нет доступа к этой команде.")
	}

	rows, err := database.DB.Query(context.Background(), `
        SELECT id, client_name, appointment_time FROM appointments
        ORDER BY appointment_time ASC
    `)
	if err != nil {
		return c.Send("⚠️ Ошибка при получении записей")
	}
	defer rows.Close()

	var result string
	for rows.Next() {
		var id int
		var clientName string
		var appointmentTime time.Time
		err := rows.Scan(&id, &clientName, &appointmentTime)
		if err != nil {
			return c.Send("⚠️ Ошибка при обработке данных")
		}
		result += fmt.Sprintf("📅 %s | %s\n", appointmentTime.Format("02.01.2006 15:04"), clientName)
	}

	if result == "" {
		return c.Send("📭 Нет записей на данный момент")
	}

	return c.Send(result)
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
	slot, exists := tempStorage[userID]
	if !exists {
		return c.Send("Ошибка: выберите слот перед выбором услуги.")
	}

	timeStr := slot.Time[11:]

	msg := fmt.Sprintf(
		"📅 Дата: %s\n💆‍♀️ Услуга: %s\n⏳ Время работы мастера: %s\n💰 Цена: %d руб\n\n"+
			"Подтвердите свою запись к своему любимому мастеру Зухре 😊",
		timeStr,
		name,
		formatDuration(duration),
		price,
	)

	btnYes := telebot.InlineButton{Text: "✅ Да", Data: fmt.Sprintf("confirm_booking|%d|%s", serviceID, timeStr)}
	btnNo := telebot.InlineButton{Text: "❌ Нет", Data: "cancel_booking"}

	return c.Send(msg, &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{{btnYes}, {btnNo}},
	})
}

func ConfirmBookingHandler(c telebot.Context) error {
	data := c.Data()
	parts := strings.Split(data, "|")
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

	var slotID int
	err := database.DB.QueryRow(
		context.Background(),
		"SELECT id FROM available_slots WHERE date = $1 AND time = $2 AND is_active = TRUE AND NOT EXISTS (SELECT 1 FROM bookings WHERE slot_id = available_slots.id)",
		formattedDate,
		formattedTime,
	).Scan(&slotID)
	if err != nil {
		return c.Send("К сожалению, этот слот уже забронирован. Выберите другой.")
	}

	_, err = database.DB.Exec(
		context.Background(),
		"INSERT INTO bookings (client_id, slot_id, service_id) VALUES ((SELECT id FROM clients WHERE telegram_id = $1), $2, $3)",
		userID, slotID, serviceID,
	)
	if err != nil {
		return c.Send("Ошибка при создании бронирования. Попробуйте снова.")
	}

	return c.Send("✅ Ваша запись отправлена на подтверждение мастеру Зухре. Ожидайте, вам придет уведомление.", &telebot.ReplyMarkup{ReplyKeyboard: MainMenu.ReplyKeyboard})
}

func MessageHandler(c telebot.Context) error {
	userID := c.Sender().ID

	switch userState[userID] {
	case StateAwaitingName:
		return processNameInput(c) // Обработка имени

	default:
		switch c.Text() {
		case "➕ Записаться к Зухре":
			return showDatePicker(c)
		default:
			return c.Send("Я не понял команду. Попробуйте снова.", &telebot.ReplyMarkup{ReplyKeyboard: MainMenu.ReplyKeyboard})
		}
	}
}
