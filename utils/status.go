package utils

// Функция для перевода статусов с эмодзи
func FormatStatus(status string) string {
	switch status {
	case "completed":
		return "✅ Завершено"
	case "cancelled":
		return "❌ Отменено"
	case "pending":
		return "⏳ В ожидании"
	default:
		return "❔ Неизвестно"
	}
}
