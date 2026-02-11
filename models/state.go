package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	StateAwaitingName                = "awaiting_name"
	StateAwaitingContact             = "awaiting_contact"
	StateNone                        = "none"
	StateAwaitingServiceName         = "awaiting_service_name"
	StateAwaitingServiceNameEdit     = "awaiting_service_name_edit"
	StateAwaitingServicePrice        = "awaiting_service_price"
	StateAwaitingServicePriceEdit    = "awaiting_service_price_edit"
	StateAwaitingServiceDuration     = "state_awaiting_service_duration"
	StateAwaitingServiceDurationEdit = "state_awaiting_service_duration_edit"
)

var TempStorage = make(map[int64]SelectedSlot)
var UserState = make(map[int64]string)
var RegistrationStorage = make(map[int64]RegistrationState)
var TempServiceData = make(map[int64]TempService)

type TempService struct {
	Name     string
	Price    int
	Duration int
	ID       int
}

type SelectedSlot struct {
	Date string
	Time string
}

type RegistrationState struct {
	Name     string
	Username string
	Phone    string
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

type Booking struct {
	ID               string    `json:"id"`
	ClientName       string    `json:"client_name"`
	ClientID         string    `json:"client_id"`
	ServiceID        string    `json:"service_id"`
	SlotID           string    `json:"slot_id"`
	DateTime         time.Time `json:"date_time"`
	Rating           *int      `json:"rating"`
	ClientTelegramID string    `json:"telegram_id"`
	ServiceName      string    `json:"service_name"`
	ServicePrice     int       `json:"service_price"`
	Status           string    `json:"status"`
}
