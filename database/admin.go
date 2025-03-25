package database

import (
	"context"
	"log"
)

var adminCache = make(map[string]bool)

func IsAdmin(telegramID string) bool {
	// Если уже есть в кэше — возвращаем
	if isAdmin, ok := adminCache[telegramID]; ok {
		return isAdmin
	}
	clientID, _ := GetUserIDByTelegram(telegramID)

	// Запрашиваем в базе
	var count int
	err := DB.QueryRow(context.Background(), "SELECT COUNT(*) FROM admins WHERE client_id = $1", clientID).Scan(&count)
	if err != nil {
		log.Println("Ошибка при проверке админа:", err)
		return false
	}

	// Кладём в кэш
	adminCache[telegramID] = count > 0
	return count > 0
}
