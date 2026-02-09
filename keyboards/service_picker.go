package keyboards

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aamarsel/browz/database"
	"gopkg.in/telebot.v3"
)

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
