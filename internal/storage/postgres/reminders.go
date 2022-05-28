package postgres

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lus/jacques/internal/reminder"
	"time"
)

// ReminderRepository implements the reminder.Repository API using PostgreSQL
type ReminderRepository struct {
	db *pgxpool.Pool
}

var _ reminder.Repository = (*ReminderRepository)(nil)

func (repo *ReminderRepository) GetByID(ctx context.Context, id uuid.UUID) (*reminder.Reminder, error) {
	row := repo.db.QueryRow(ctx, "SELECT * FROM reminders WHERE reminder_id = $1", id)
	obj, err := repo.rowToReminder(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}

func (repo *ReminderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*reminder.Reminder, error) {
	rows, err := repo.db.Query(ctx, "SELECT * FROM reminders WHERE user_id = $1", userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var reminders []*reminder.Reminder
	for rows.Next() {
		obj, err := repo.rowToReminder(rows)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, obj)
	}
	return reminders, nil
}

func (repo *ReminderRepository) GetNext(ctx context.Context) (*reminder.Reminder, error) {
	row := repo.db.QueryRow(ctx, "SELECT * FROM reminders ORDER BY fires_at")
	obj, err := repo.rowToReminder(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}

func (repo *ReminderRepository) Create(ctx context.Context, create *reminder.Create) (*reminder.Reminder, error) {
	obj := &reminder.Reminder{
		ID:          uuid.New(),
		UserID:      create.UserID,
		ChannelID:   create.ChannelID,
		Description: create.Description,
		FiresAt:     time.Now().Add(create.Delta).Unix(),
		CreatedAt:   time.Now().Unix(),
	}

	query := "INSERT INTO reminders VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := repo.db.Exec(ctx, query, obj.ID, obj.UserID, obj.ChannelID, obj.Description, obj.FiresAt, obj.CreatedAt)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (repo *ReminderRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.Exec(ctx, "DELETE FROM reminders WHERE reminder_id = $1", id)
	return err
}

func (repo *ReminderRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := repo.db.Exec(ctx, "DELETE FROM reminders WHERE user_id = $1", userID)
	return err
}

func (repo *ReminderRepository) rowToReminder(row pgx.Row) (*reminder.Reminder, error) {
	obj := new(reminder.Reminder)
	err := row.Scan(
		&obj.ID,
		&obj.UserID,
		&obj.ChannelID,
		&obj.Description,
		&obj.FiresAt,
		&obj.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
