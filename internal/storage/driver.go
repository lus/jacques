package storage

import (
	"context"
	"github.com/lus/jacques/internal/reminder"
)

// Driver defines the storage driver API
type Driver interface {
	Initialize(ctx context.Context) error
	Reminders() reminder.Repository
	Close() error
}
