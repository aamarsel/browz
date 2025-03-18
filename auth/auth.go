package auth

import (
	"context"

	"github.com/aamarsel/browz/database"
)

func IsAdmin(userID int64) bool {
	var exists bool
	err := database.DB.QueryRow(context.Background(), `
        SELECT EXISTS(SELECT 1 FROM admins WHERE telegram_id = $1)
    `, userID).Scan(&exists)

	return err == nil && exists
}
