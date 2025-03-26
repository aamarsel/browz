package database

import (
	"context"
	"log"

	"github.com/aamarsel/browz/models"
)

func GetClientByID(clientID int64) (*models.Client, error) {
	query := `SELECT id, name, phone, telegram_id FROM clients WHERE telegram_id = $1`
	var client models.Client

	err := DB.QueryRow(context.Background(), query, clientID).Scan(
		&client.ID, &client.Name, &client.Phone, &client.TelegramID,
	)
	if err != nil {
		log.Println("Ошибка при получении данных клиента:", err)
		return nil, err
	}

	return &client, nil
}
