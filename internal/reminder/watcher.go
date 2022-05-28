package reminder

import (
	"context"
	"github.com/rs/zerolog/log"
	"time"
)

type Callback func(fired *Reminder)

// Watcher represents the service that watches reminders and executes fire events
type Watcher struct {
	Repo Repository

	callbacks []Callback

	next       *Reminder
	minFetchAt time.Time

	running bool
	stop    chan struct{}
}

// Start starts watching for reminders to fire
func (watcher *Watcher) Start() {
	if watcher.running {
		return
	}
	watcher.stop = make(chan struct{})
	go func() {
		for {
			select {
			case <-watcher.stop:
				return
			default:
				if watcher.next == nil && time.Until(watcher.minFetchAt) <= 0 {
					next, err := watcher.Repo.GetNext(context.Background())
					if err != nil {
						log.Error().Err(err).Msg("could not fetch next reminder")
						watcher.minFetchAt = time.Now().Add(30 * time.Second)
					} else if next == nil {
						watcher.minFetchAt = time.Now().Add(time.Minute)
					}
					watcher.next = next
				} else if watcher.next != nil && time.Until(time.Unix(watcher.next.FiresAt, 0)) <= 0 {
					for _, callback := range watcher.callbacks {
						callback(watcher.next)
					}
					if err := watcher.Repo.DeleteByID(context.Background(), watcher.next.ID); err != nil {
						log.Error().Err(err).Msg("could not delete fired reminder")
					}
					watcher.Reset()
				}
			}
		}
	}()
	watcher.running = true
}

// Stop stops watching for reminders to fire
func (watcher *Watcher) Stop() {
	if !watcher.running {
		return
	}
	close(watcher.stop)
	watcher.running = false
}

// Reset resets the cached next reminder
func (watcher *Watcher) Reset() {
	watcher.next = nil
	watcher.minFetchAt = time.Unix(0, 0)
}

// Subscribe registers a callback for reminder fire events
func (watcher *Watcher) Subscribe(callback Callback) {
	watcher.callbacks = append(watcher.callbacks, callback)
}
