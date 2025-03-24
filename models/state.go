package models

const (
	StateAwaitingName    = "awaiting_name"
	StateAwaitingContact = "awaiting_contact"
	StateNone            = "none"
)

var TempStorage = make(map[int64]SelectedSlot)
var UserState = make(map[int64]string)
var RegistrationStorage = make(map[int64]RegistrationState)

type SelectedSlot struct {
	Date string
	Time string
}

type RegistrationState struct {
	Name  string
	Phone string
}
