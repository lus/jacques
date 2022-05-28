package reminder

import (
	"context"
	"github.com/google/uuid"
	"time"
)

// Repository defines the reminder repository API
type Repository interface {
	// GetByID retrieves a reminder by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*Reminder, error)

	// GetByUserID retrieves all reminders of a user
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Reminder, error)

	// GetNext retrieves the reminder that fires next
	GetNext(ctx context.Context) (*Reminder, error)

	// Create creates a new reminder
	Create(ctx context.Context, create *Create) (*Reminder, error)

	// DeleteByID deletes a reminder by its ID
	DeleteByID(ctx context.Context, id uuid.UUID) error

	// DeleteByUserID deletes all reminders of a user
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

// Create represents the reminder creation payload
type Create struct {
	UserID      int64
	ChannelID   int64
	Description string
	Delta       time.Duration
}
