package keyboards

import (
	"github.com/aamarsel/browz/database"
	"gopkg.in/telebot.v3"
)

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
