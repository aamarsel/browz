package keyboards

import (
	"fmt"
	"time"

	"gopkg.in/telebot.v3"
)

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
