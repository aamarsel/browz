package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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
	userID := c.Sender().ID

	tempStorage[userID] = SelectedSlot{
		Date: data[:10],
		Time: data[11:],
	}

	// Показываем список услуг
	return c.Send("Выберите услугу:", &telebot.ReplyMarkup{
		InlineKeyboard: GetServicesButtons(),
	})
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

func GetServicesButtons() [][]telebot.InlineButton {
	rows, err := database.DB.Query(context.Background(), "SELECT id, name, price, duration FROM services")
	if err != nil {
		log.Println("Ошибка при получении списка услуг:", err)
		return nil
	}
	defer rows.Close()

	var buttons [][]telebot.InlineButton
	for rows.Next() {
		var id int
		var name string
		var price int
		var duration time.Duration
		err := rows.Scan(&id, &name, &price, &duration)
		if err != nil {
			log.Println("Ошибка при обработке строки услуги:", err)
			continue
		}

		btn := telebot.InlineButton{
			Text: fmt.Sprintf("%s, %s, %d руб", name, formatDuration(duration), price),
			Data: fmt.Sprintf("pick_service:%d", id),
		}
		buttons = append(buttons, []telebot.InlineButton{btn})
	}
	return buttons
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	if minutes == 0 {
		return fmt.Sprintf("%d ч", hours)
	}
	return fmt.Sprintf("%d ч %d мин", hours, minutes)
}

// Главное меню для клиентов
var MainMenu = &telebot.ReplyMarkup{}

var btnMyBookings = MainMenu.Text("📅 Мои бронирования")
var btnNewBooking = MainMenu.Text("➕ Записаться к Зухре")

func InitKeyboards() {
	MainMenu.Reply(
		MainMenu.Row(btnMyBookings),
		MainMenu.Row(btnNewBooking),
	)
}
