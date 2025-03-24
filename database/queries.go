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
