package reminder

import (
	"context"
	"github.com/google/uuid"
)

type watcherResettingRepository struct {
	wrapping Repository
	watcher  *Watcher
}

var _ Repository = (*watcherResettingRepository)(nil)

func (repo *watcherResettingRepository) GetByID(ctx context.Context, id uuid.UUID) (*Reminder, error) {
	return repo.wrapping.GetByID(ctx, id)
}

func (repo *watcherResettingRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Reminder, error) {
	return repo.wrapping.GetByUserID(ctx, userID)
}

func (repo *watcherResettingRepository) GetNext(ctx context.Context) (*Reminder, error) {
	return repo.wrapping.GetNext(ctx)
}

func (repo *watcherResettingRepository) Create(ctx context.Context, create *Create) (*Reminder, error) {
	obj, err := repo.wrapping.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	repo.watcher.Reset()
	return obj, nil
}

func (repo *watcherResettingRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	err := repo.wrapping.DeleteByID(ctx, id)
	if err != nil {
		return err
	}
	repo.watcher.Reset()
	return nil
}

func (repo *watcherResettingRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	err := repo.wrapping.DeleteByUserID(ctx, userID)
	if err != nil {
		return err
	}
	repo.watcher.Reset()
	return nil
}
