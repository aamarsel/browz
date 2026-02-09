package keyboards

import (
	"fmt"

	"github.com/aamarsel/browz/database"
	"gopkg.in/telebot.v3"
)

// Главное меню
var MainMenu = &telebot.ReplyMarkup{}

var btnMyBookings = MainMenu.Text("📅 Мои бронирования")
var btnNewBooking = MainMenu.Text("➕ Записаться к Зухре")
var btnFutureBookings = MainMenu.Text("📅 Будущие записи")
var btnPendingBookings = MainMenu.Text("⏳ Неподтвержденные записи")
var btnOldBookings = MainMenu.Text("📜 Прошлые записи")
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
			menu.Row(btnOldBookings),
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
