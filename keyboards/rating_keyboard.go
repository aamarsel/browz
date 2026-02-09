package keyboards

import (
	"fmt"

	"github.com/google/uuid"
	"gopkg.in/telebot.v3"
)

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
