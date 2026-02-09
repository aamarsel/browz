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
