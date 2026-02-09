package keyboards

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aamarsel/browz/database"
	"github.com/google/uuid"
	"gopkg.in/telebot.v3"
)

// Главное меню
var MainMenu = &telebot.ReplyMarkup{}

var btnMyBookings = MainMenu.Text("📅 Мои бронирования")
var btnNewBooking = MainMenu.Text("➕ Записаться к Зухре")
var btnFutureBookings = MainMenu.Text("📅 Будущие записи")
var btnPendingBookings = MainMenu.Text("⏳ Неподтвержденные записи")
var btnEditServices = MainMenu.Text("✏️ Редактировать услуги")
var btnNewService = MainMenu.Text("➕ Добавить услугу")

func GetMainMenu(isAdmin bool) *telebot.ReplyMarkup {
	menu := &telebot.ReplyMarkup{}

	// Кнопки для всех пользователей
	menu.Reply(
		menu.Row(btnMyBookings),
		menu.Row(btnNewBooking),
	)

	// Если админ — добавляем доп. кнопки
	if isAdmin {
		menu.Reply(
			menu.Row(btnFutureBookings),
			menu.Row(btnPendingBookings),
			menu.Row(btnMyBookings),
			menu.Row(btnNewBooking),
			menu.Row(btnEditServices),
			menu.Row(btnNewService),
		)
	}

	return menu
}

func SendMainMenu(c telebot.Context, text string) error {
	isAdmin := database.IsAdmin(fmt.Sprint(c.Sender().ID)) // Проверяем один раз
	menu := GetMainMenu(isAdmin)
	return c.Send(text, menu)
}

// Показ календаря с выбором даты
func ShowDatePicker(c telebot.Context, page int) error {
	today := time.Now()
	btns := &telebot.ReplyMarkup{}

	var rows []telebot.Row
	for i := 0; i < 7; i++ {
		date := today.AddDate(0, 0, i+page*7)
		btn := btns.Data(date.Format("02.01.2006"), "pick_date", date.Format("2006-01-02"))
		rows = append(rows, btns.Row(btn))
	}
	rows = append(rows, btns.Row(btns.Data(">>", "next_date_page", fmt.Sprintf("%d", page+1))))
	if page > 0 {
		rows = append(rows, btns.Row(btns.Data("<<", "prev_date_page", fmt.Sprintf("%d", page-1))))
	}

	btns.Inline(rows...)
	return c.Send("Выберите дату посещения:", btns)
}

// Показ доступных слотов
func ShowTimeSlots(c telebot.Context, date string) error {
	slots, err := database.GetAvailableSlots(date)
	if err != nil {
		return c.Send("Ошибка при загрузке доступных слотов")
	}
	if len(slots) == 0 {
		c.Send("На этот день все слоты заняты. Выберите другой день")
		return ShowDatePicker(c, 0)
	}

	btns := &telebot.ReplyMarkup{ResizeKeyboard: true}
	var rows []telebot.Row
	for _, slot := range slots {
		btn := btns.Data(slot, "pick_slot", date+" "+slot)
		rows = append(rows, telebot.Row{btn})
	}
	btns.Inline(rows...)
	return c.Send("Выберите время:", btns)
}

func RatingKeyboard(bookingID uuid.UUID) *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}

	var rows []telebot.Row
	for i := 1; i <= 5; i++ {
		btn := markup.Data(
			fmt.Sprintf("%d ⭐", i),
			"rate_service",
			fmt.Sprintf("%s|%d", bookingID, i),
		)
		rows = append(rows, markup.Row(btn))
	}

	markup.Inline(rows...)
	return markup
}

// Получение кнопок с услугами
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
			Text: fmt.Sprintf("%s, %d руб", name, price),
			Data: fmt.Sprintf("pick_service:%d", id),
		}
		buttons = append(buttons, []telebot.InlineButton{btn})
	}
	return buttons
}

// Вспомогательная функция форматирования времени
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	if minutes == 0 {
		return fmt.Sprintf("%d ч", hours)
	}
	return fmt.Sprintf("%d ч %d мин", hours, minutes)
}
