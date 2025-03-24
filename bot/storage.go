package bot

const (
	StateAwaitingName    = "awaiting_name"
	StateAwaitingContact = "awaiting_contact"
)

var tempStorage = make(map[int64]SelectedSlot)
var userState = make(map[int64]string)
var registrationStorage = make(map[int64]RegistrationState)

type SelectedSlot struct {
	Date string
	Time string
}

type RegistrationState struct {
	Name  string
	Phone string
}
