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
	ID          string    `json:"id"`
	DateTime    time.Time `json:"date_time"`
	ServiceName string    `json:"service_name"`
	Status      string    `json:"status"`
}
