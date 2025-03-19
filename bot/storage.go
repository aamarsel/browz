package bot

var tempStorage = make(map[int64]SelectedSlot)

type SelectedSlot struct {
	Date string
	Time string
}
