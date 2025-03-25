package database

import "time"

type Appointment struct {
	ID         int
	ClientName string
	Phone      string
	Date       time.Time
	Time       string
	Status     string
}

type Booking struct {
	ID               string    `json:"id"`
	ClientName       string    `json:"client_name"`
	ClientID         string    `json:"client_id"`
	ServiceID        string    `json:"service_id"`
	SlotID           string    `json:"slot_id"`
	DateTime         time.Time `json:"date_time"`
	ClientTelegramID string    `json:"telegram_id"`
	ServiceName      string    `json:"service_name"`
	Status           string    `json:"status"`
}
