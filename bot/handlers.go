package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/aamarsel/browz/auth"
	"github.com/aamarsel/browz/database"
)

func StartHandler(c telebot.Context) error {
	return c.Send("Привет! Я бот для записи на брови. Введите /book для записи.")
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
	} else {
		log.Println("Ошибка! Неизвестный callback")
	}

	return nil
}

func ProcessBooking(c telebot.Context) error {
	data := strings.Fields(c.Text())
	if len(data) < 4 {
		return c.Send("Ошибка! Введите данные в формате: Имя Телефон Дата Время (пример: Анна +79998887766 2025-03-20 14:00)")
	}

	name, phone, date, time := data[0], data[1], data[2], data[3]

	_, err := database.DB.Exec(context.Background(),
		"INSERT INTO appointments (client_name, phone, appointment_date, appointment_time) VALUES ($1, $2, $3, $4)",
		name, phone, date, time)

	if err != nil {
		log.Println("Ошибка при записи в БД:", err)
		return c.Send("Ошибка при записи! Попробуйте еще раз.")
	}

	return c.Send(fmt.Sprintf("✅ %s, вы записаны на %s в %s!", name, date, time))
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
