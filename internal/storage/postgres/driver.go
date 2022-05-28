package postgres

import (
	"context"
	"embed"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lus/jacques/internal/reminder"
	"github.com/lus/jacques/internal/storage"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Driver represents the PostgreSQL storage driver
type Driver struct {
	dsn       string
	db        *pgxpool.Pool
	reminders *ReminderRepository
}

var _ storage.Driver = (*Driver)(nil)

func New(dsn string) *Driver {
	return &Driver{
		dsn: dsn,
	}
}

func (driver *Driver) Initialize(ctx context.Context) error {
	// Perform SQL migrations
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}
	migrator, err := migrate.NewWithSourceInstance("iofs", source, driver.dsn)
	if err != nil {
		return err
	}
	defer migrator.Close()
	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	// Initialize the database connection pool
	pool, err := pgxpool.Connect(ctx, driver.dsn)
	if err != nil {
		return err
	}
	driver.db = pool

	// Instantiate the repositories
	driver.reminders = &ReminderRepository{db: pool}

	return nil
}

func (driver *Driver) Reminders() reminder.Repository {
	return driver.reminders
}

func (driver *Driver) Close() error {
	driver.reminders = nil

	driver.db.Close()
	driver.db = nil
	return nil
}
