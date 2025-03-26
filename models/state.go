package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	StateAwaitingName            = "awaiting_name"
	StateAwaitingContact         = "awaiting_contact"
	StateNone                    = "none"
	StateAwaitingServiceName     = "awaiting_service_name"
	StateAwaitingServicePrice    = "awaiting_service_price"
	StateAwaitingServiceDuration = "state_awaiting_service_duration"
)

var TempStorage = make(map[int64]SelectedSlot)
var UserState = make(map[int64]string)
var RegistrationStorage = make(map[int64]RegistrationState)
var TempServiceData = make(map[int64]TempService)

type TempService struct {
	Name     string
	Price    int
	Duration int
}

type SelectedSlot struct {
	Date string
	Time string
}

type RegistrationState struct {
	Name  string
	Phone string
}

type Service struct {
	ID       int
	Name     string
	Price    int
	Duration time.Duration
}

type Client struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	TelegramID int64     `json:"telegram_id"`
}
