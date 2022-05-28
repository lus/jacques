package reminder

import (
	"context"
	"github.com/google/uuid"
)

type WatcherResettingRepository struct {
	Wrapping Repository
	Watcher  *Watcher
}

var _ Repository = (*WatcherResettingRepository)(nil)

func (repo *WatcherResettingRepository) GetByID(ctx context.Context, id uuid.UUID) (*Reminder, error) {
	return repo.Wrapping.GetByID(ctx, id)
}

func (repo *WatcherResettingRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Reminder, error) {
	return repo.Wrapping.GetByUserID(ctx, userID)
}

func (repo *WatcherResettingRepository) GetNext(ctx context.Context) (*Reminder, error) {
	return repo.Wrapping.GetNext(ctx)
}

func (repo *WatcherResettingRepository) Create(ctx context.Context, create *Create) (*Reminder, error) {
	obj, err := repo.Wrapping.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	repo.Watcher.Reset()
	return obj, nil
}

func (repo *WatcherResettingRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	err := repo.Wrapping.DeleteByID(ctx, id)
	if err != nil {
		return err
	}
	repo.Watcher.Reset()
	return nil
}

func (repo *WatcherResettingRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	err := repo.Wrapping.DeleteByUserID(ctx, userID)
	if err != nil {
		return err
	}
	repo.Watcher.Reset()
	return nil
}
