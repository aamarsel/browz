package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aamarsel/browz/database"
	"github.com/aamarsel/browz/keyboards"
	"github.com/google/uuid"
	"gopkg.in/telebot.v3"
)

func CompleteOldBookings(bot *telebot.Bot) {
	now := time.Now()

	// Запрос на выборку всех записей, которые уже должны были закончиться
	query := `
		SELECT b.id, b.client_id, b.service_id, s.name, sl.date, sl.time, s.duration, c.telegram_id, c.name
		FROM bookings b
		JOIN services s ON b.service_id = s.id
		JOIN available_slots sl ON b.slot_id = sl.id
		JOIN clients c ON b.client_id = c.id
		WHERE b.status = 'accepted'
	`

	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		log.Println("Ошибка при выборке старых записей:", err)
		return
	}
	defer rows.Close()

	// Буферизуем результаты в слайс, чтобы закрыть соединение перед транзакцией
	var bookings []struct {
		BookingID        uuid.UUID
		ClientID         uuid.UUID
		ServiceID        int64
		ServiceName      string
		Date             time.Time
		TimeOfDay        time.Time
		Duration         time.Duration
		ClientTelegramID int64
		ClientName       string
	}

	for rows.Next() {
		var b struct {
			BookingID        uuid.UUID
			ClientID         uuid.UUID
			ServiceID        int64
			ServiceName      string
			Date             time.Time
			TimeOfDay        time.Time
			Duration         time.Duration
			ClientTelegramID int64
			ClientName       string
		}

		err := rows.Scan(&b.BookingID, &b.ClientID, &b.ServiceID, &b.ServiceName, &b.Date, &b.TimeOfDay, &b.Duration, &b.ClientTelegramID, &b.ClientName)
		if err != nil {
			log.Println("Ошибка при сканировании записи:", err)
			continue
		}
		bookings = append(bookings, b)
	}

	// Закрываем соединение перед обновлением статусов
	rows.Close()

	// Теперь обрабатываем записи
	for _, b := range bookings {
		// Вычисляем время окончания услуги
		endTime := time.Date(b.Date.Year(), b.Date.Month(), b.Date.Day(), b.TimeOfDay.Hour(), b.TimeOfDay.Minute(), b.TimeOfDay.Second(), 0, time.Local).Add(b.Duration)
		log.Println(endTime, now.After(endTime), now)

		if now.After(endTime) {
			// Открываем новую транзакцию
			tx, err := database.DB.Begin(context.Background())
			if err != nil {
				log.Println("Ошибка при создании транзакции:", err)
				continue
			}

			// Обновляем статус на "completed"
			_, err = tx.Exec(context.Background(),
				`UPDATE bookings SET status = 'completed' WHERE id = $1`, b.BookingID)
			if err != nil {
				log.Println("Ошибка при обновлении статуса записи:", err)
				tx.Rollback(context.Background()) // Откат транзакции
				continue
			}

			// Фиксируем изменения
			err = tx.Commit(context.Background())
			if err != nil {
				log.Println("Ошибка при фиксации транзакции:", err)
				continue
			}

			// Отправляем сообщение клиенту
			message := fmt.Sprintf(
				"🌸 *%s*, спасибо, что воспользовались услугами Зухры!\n\n"+
					"Мы будем рады видеть вас снова. Вы можете записаться заранее, чтобы выбрать удобное время. 💕",
				b.ClientName,
			)

			recipient := &telebot.User{ID: b.ClientTelegramID}
			bot.Send(recipient, message, telebot.ModeMarkdown)
			bot.Send(recipient, "Пожалуйста, оцените качество услуги:", keyboards.RatingKeyboard(b.BookingID))
		}
	}
}
