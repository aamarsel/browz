package database

import (
	"context"
	"log"
	"time"
)

// Получение записей по статусу
func GetBookingsByStatus(status string) ([]Booking, error) {
	query := `
	SELECT b.id, c.name, s.name, sl.date, sl.time, b.status
		FROM bookings b
		JOIN clients c ON b.client_id = c.id
		JOIN services s ON b.service_id = s.id
		JOIN available_slots sl ON b.slot_id = sl.id
		WHERE b.status = $1
		ORDER BY sl.date, sl.time;
	`
	rows, err := DB.Query(context.Background(), query, status)
	if err != nil {
		log.Println("Ошибка при получении записей:", err)
		return nil, err
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var b Booking
		var date time.Time
		var timeStr string
		if err := rows.Scan(&b.ID, &b.ClientName, &b.ServiceName, &date, &timeStr, &b.Status); err != nil {
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

// Обновление статуса записи
func UpdateBookingStatus(bookingID, newStatus string) (*Booking, error) {
	query := `
		UPDATE bookings 
		SET status = $1 
		WHERE id = $2 
		RETURNING id, client_id, service_id, slot_id, status
	`
	var booking Booking
	err := DB.QueryRow(context.Background(), query, newStatus, bookingID).Scan(
		&booking.ID, &booking.ClientID, &booking.ServiceID, &booking.SlotID, &booking.Status,
	)
	if err != nil {
		log.Println("Ошибка при обновлении статуса записи:", err)
		return nil, err
	}

	var date time.Time
	var timeOfDay time.Time // Меняем с time.Duration на time.Time

	// Получаем дополнительные данные
	err = DB.QueryRow(context.Background(),
		`SELECT s.name, sl.date, sl.time, c.telegram_id 
		 FROM services s
		 JOIN available_slots sl ON sl.id = $1
		 JOIN clients c ON c.id = $2
		 WHERE s.id = $3`,
		booking.SlotID, booking.ClientID, booking.ServiceID,
	).Scan(&booking.ServiceName, &date, &timeOfDay, &booking.ClientTelegramID)
	if err != nil {
		log.Println("Ошибка при получении деталей бронирования:", err)
		return nil, err
	}

	booking.DateTime = time.Date(
		date.Year(), date.Month(), date.Day(),
		timeOfDay.Hour(), timeOfDay.Minute(), timeOfDay.Second(), 0,
		date.Location(),
	)

	return &booking, nil
}
