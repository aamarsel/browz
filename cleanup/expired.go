package cleanup

import (
	"context"
	"log"

	"github.com/aamarsel/browz/database"
)

func ClearExpired() {
	res, err := database.DB.Exec(
		context.Background(),
		`
			UPDATE bookings b
			SET status = 'cancelled'
			FROM available_slots s
			WHERE b.slot_id = s.id
			AND b.status = 'pending'
			AND (s.date < CURRENT_DATE OR (s.date = CURRENT_DATE AND s.time < CURRENT_TIME));
		`,
	)
	if err != nil {
		log.Fatal(err)
	}

	rows := res.RowsAffected()
	log.Printf("expired %d bookings", rows)
}
