package database

import (
	"context"
	"fmt"
	"log"

	"github.com/aamarsel/browz/models"
)

func GetAllServices() ([]models.Service, error) {
	rows, err := DB.Query(context.Background(), "SELECT id, name, price, duration FROM services ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var s models.Service
		if err := rows.Scan(&s.ID, &s.Name, &s.Price, &s.Duration); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

func DeleteService(serviceID int) error {
	_, err := DB.Exec(context.Background(), "DELETE FROM services WHERE id = $1", serviceID)
	return err
}

func AddService(s models.TempService) error {
	durationStr := fmt.Sprintf("%d minutes", s.Duration)

	_, err := DB.Exec(context.Background(),
		"INSERT INTO services (name, price, duration) VALUES ($1, $2, $3)",
		s.Name, s.Price, durationStr)
	return err
}

func GetServiceByID(serviceID int) (*models.Service, error) {
	query := `SELECT id, name, price, duration FROM services WHERE id = $1`
	var service models.Service

	err := DB.QueryRow(context.Background(), query, serviceID).Scan(
		&service.ID, &service.Name, &service.Price, &service.Duration,
	)
	if err != nil {
		log.Println("Ошибка при получении данных об услуге:", err)
		return nil, err
	}

	return &service, nil
}
