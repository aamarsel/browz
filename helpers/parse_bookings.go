package helpers

import (
	"log"
	"time"

	"github.com/aamarsel/browz/models"
	"github.com/jackc/pgx/v5"
)

func ParseBookingsFromRows(rows pgx.Rows) ([]models.Booking, error) {
	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		var date time.Time
		var timeStr string

		err := rows.Scan(&b.ID, &b.ClientName, &b.ServiceName, &date, &timeStr, &b.Status, &b.Rating)
		if err != nil {
			log.Println("Ошибка при обработке записи:", err)
			continue
		}

		// Парсим время и объединяем с датой
		parsedTime, _ := time.Parse("15:04:05", timeStr)
		b.DateTime = time.Date(date.Year(), date.Month(), date.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.UTC)

		bookings = append(bookings, b)
	}
	return bookings, nil
}
