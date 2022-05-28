package reminder

import (
	"github.com/google/uuid"
)

// Reminder represents a reminder created by a user
type Reminder struct {
	ID          uuid.UUID
	UserID      int64
	ChannelID   int64
	Description string
	FiresAt     int64
	CreatedAt   int64
}
