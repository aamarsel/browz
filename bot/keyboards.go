package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"gopkg.in/telebot.v3"

	"github.com/aamarsel/browz/database"
)

var daysOfWeek = map[time.Weekday]string{
	time.Monday:    "понедельник",
	time.Tuesday:   "вторник",
	time.Wednesday: "среда",
	time.Thursday:  "четверг",
	time.Friday:    "пятница",
	time.Saturday:  "суббота",
	time.Sunday:    "воскресенье",
}

func showDatePicker(c telebot.Context) error {
	today := time.Now()
	btns := &telebot.ReplyMarkup{}

	var rows []telebot.Row
	for i := 0; i < 7; i++ {
		date := today.AddDate(0, 0, i)
		btn := btns.Data(date.Format("02.01.2006"), "pick_date", date.Format("2006-01-02"))
		rows = append(rows, btns.Row(btn))
	}

	btns.Inline(rows...)

	return c.Send("Выберите дату посещения:", btns)
}

func DatePickerHandler(c telebot.Context) error {
	parts := strings.Split(c.Data(), "|")
	if len(parts) < 2 {
		return c.Send("Ошибка: неверный формат данных.")
	}

	selectedDate := parts[1]
	return showTimeSlots(c, selectedDate)
}

func showTimeSlots(c telebot.Context, date string) error {
	slots, err := getAvailableSlots(date)
	if err != nil {
		return c.Send("Ошибка при загрузке доступных слотов")
	}
	if len(slots) == 0 {
		c.Send("На этот день все слоты заняты. Выберите другой день")
		return showDatePicker(c)
	}

	btns := &telebot.ReplyMarkup{ResizeKeyboard: true}
	var rows []telebot.Row
	for _, slot := range slots {
		btn := btns.Data(slot, "pick_slot", date+" "+slot)
		rows = append(rows, telebot.Row{btn})
	}
	btns.Inline(rows...)

	return c.Send("Выберите время, удобное для записи. Все слоты внизу - свободные:", btns)
}

func SlotPickerHandler(c telebot.Context) error {
	data := c.Data()
	return createBooking(c, data)
}

func createBooking(c telebot.Context, data string) error {
	log.Println("Data:", data)
	userID := c.Sender().ID

	// Разбираем дату и время из data
	date := data[11:21] // "2025-03-20"
	time := data[34:39] // "08:00"
	log.Println("LOG:", date, time)

	// Получаем ID клиента по его Telegram ID
	var clientID uuid.UUID
	err := database.DB.QueryRow(
		context.Background(),
		"SELECT id FROM clients WHERE telegram_id = $1",
		userID,
	).Scan(&clientID)
	if err != nil {
		return c.Send("Ошибка: вы не зарегистрированы.")
	}

	// Получаем ID слота, который клиент хочет забронировать
	var slotID int
	err = database.DB.QueryRow(
		context.Background(),
		`SELECT id FROM available_slots 
			WHERE date = $1 
			AND time = $2 
			AND is_active = TRUE 
			AND NOT EXISTS (SELECT 1 FROM bookings WHERE slot_id = available_slots.id)`,
		date, time,
	).Scan(&slotID)
	if err != nil {
		return c.Send("Ошибка: этот слот уже забронирован или недоступен.")
	}

	// Создаём бронирование
	_, err = database.DB.Exec(
		context.Background(),
		`INSERT INTO bookings (client_id, slot_id) 
		 VALUES ($1, $2) 
		 ON CONFLICT (client_id, slot_id) DO NOTHING`,
		clientID, slotID,
	)
	if err != nil {
		return c.Send("Ошибка при записи! Попробуйте еще раз.")
	}

	return c.Send("✅ Вы успешно записаны!")
}

func getAvailableSlots(date string) ([]string, error) {
	var slots []string
	query := `
		SELECT s.date, s.time 
		FROM available_slots s
		LEFT JOIN bookings b ON s.id = b.slot_id
		WHERE s.date = $1 
		  AND s.is_active = TRUE 
		  AND b.id IS NULL  
		ORDER BY s.time;
	`

	rows, err := database.DB.Query(context.Background(), query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var slotDate time.Time
		var slotTime string

		if err := rows.Scan(&slotDate, &slotTime); err != nil {
			return nil, err
		}

		// Парсим время
		parsedTime, err := time.Parse("15:04:05", slotTime)
		if err != nil {
			log.Println("Ошибка парсинга времени:", err)
			continue
		}

		// Переводим день недели на русский
		dayOfWeek := daysOfWeek[slotDate.Weekday()]

		// Форматируем: 20.03.2025, 08:00, понедельник
		formatted := fmt.Sprintf("%02d.%02d.%d, %02d:%02d, %s",
			slotDate.Day(), slotDate.Month(), slotDate.Year(),
			parsedTime.Hour(), parsedTime.Minute(),
			dayOfWeek,
		)

		slots = append(slots, formatted)
	}
	return slots, nil
}
