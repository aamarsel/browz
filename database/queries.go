package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aamarsel/browz/utils"
)

func ClientExists(userID int64) (bool, error) {
	var exists bool
	err := DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM clients WHERE telegram_id = $1)", userID).Scan(&exists)
	return exists, err
}

func SaveClient(userID int64, name, phone string) error {
	_, err := DB.Exec(context.Background(),
		"INSERT INTO clients (telegram_id, name, phone) VALUES ($1, $2, $3)",
		userID, name, phone)
	return err
}

func FindAvailableSlot(date string, timeStr string) (int, error) {
	var slotID int
	err := DB.QueryRow(context.Background(),
		"SELECT id FROM available_slots WHERE date = $1 AND time = $2 AND is_active = TRUE",
		date, timeStr).Scan(&slotID)

	if err != nil {
		return 0, errors.New("slot not available")
	}
	return slotID, nil
}

func BookSlot(userID int64, slotID, serviceID int) error {
	_, err := DB.Exec(context.Background(),
		"INSERT INTO bookings (client_id, slot_id, service_id) VALUES ((SELECT id FROM clients WHERE telegram_id = $1), $2, $3)",
		userID, slotID, serviceID)
	return err
}

func GetAvailableSlots(date string) ([]string, error) {
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

	rows, err := DB.Query(context.Background(), query, date)
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

		parsedTime, err := time.Parse("15:04:05", slotTime)
		if err != nil {
			log.Println("Ошибка парсинга времени:", err)
			continue
		}

		dayOfWeek := utils.GetDayOfWeek(slotDate)
		formatted := fmt.Sprintf("%02d.%02d.%d, %02d:%02d, %s",
			slotDate.Day(), slotDate.Month(), slotDate.Year(),
			parsedTime.Hour(), parsedTime.Minute(),
			dayOfWeek,
		)

		slots = append(slots, formatted)
	}
	return slots, nil
}

func GetUserBookings(telegramID string) ([]Booking, error) {
	query := `
		SELECT 
			b.id, 
			s.date, 
			s.time, 
			srv.name AS service_name,
			srv.price AS service_price, 
			b.status
		FROM bookings b
		JOIN available_slots s ON b.slot_id = s.id
		JOIN services srv ON b.service_id = srv.id
		WHERE b.client_id = $1
		ORDER BY s.date, s.time;
	`

	clientID, _ := GetUserIDByTelegram(telegramID)
	rows, err := DB.Query(context.Background(), query, clientID)
	if err != nil {
		log.Println("Ошибка при получении бронирований:", err)
		return nil, err
	}
	defer rows.Close()

	var bookings []Booking

	for rows.Next() {
		var booking Booking
		var slotDate time.Time
		var slotTime string
		log.Println("found row")

		err := rows.Scan(&booking.ID, &slotDate, &slotTime, &booking.ServiceName, &booking.ServicePrice, &booking.Status)
		if err != nil {
			log.Println("Ошибка при сканировании строки бронирования:", err)
			continue
		}

		// Парсим строковое время в формат Go
		parsedTime, err := time.Parse("15:04:05", slotTime)
		if err != nil {
			log.Println("Ошибка парсинга времени:", err)
			continue
		}

		// Собираем полное время бронирования
		booking.DateTime = time.Date(
			slotDate.Year(), slotDate.Month(), slotDate.Day(),
			parsedTime.Hour(), parsedTime.Minute(), 0, 0,
			time.Local,
		)

		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func GetUserIDByTelegram(telegramID string) (string, error) {
	var clientID string
	query := "SELECT id FROM clients WHERE telegram_id = $1"
	err := DB.QueryRow(context.Background(), query, telegramID).Scan(&clientID)
	return clientID, err
}

// CancelBooking отменяет бронирование, если оно ещё актуально
func CancelBooking(telegramID string, bookingID string) error {
	clientID, _ := GetUserIDByTelegram(telegramID)
	log.Println(bookingID, clientID)

	query := `
		UPDATE bookings 
		SET status = 'cancelled' 
		WHERE id = $1 AND client_id = $2 AND status NOT IN ('cancelled', 'completed')
	`
	res, err := DB.Exec(context.Background(), query, bookingID, clientID)
	if err != nil {
		return err
	}

	// Проверяем, было ли обновлено хотя бы 1 бронирование
	if res.RowsAffected() == 0 {
		return fmt.Errorf("не удалось отменить запись (возможно, она уже отменена или завершена)")
	}

	return nil
}
