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

	running bool
	stop    chan struct{}
}

// Start starts watching for reminders to fire
func (watcher *Watcher) Start() {
	if watcher.running {
		return
	}
	watcher.stop = make(chan struct{})

	// Retrieve the next reminder
	next, err := watcher.Repo.GetNext(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("could not fetch next reminder")
		watcher.Stop()
		return
	}
	if next == nil {
		watcher.Stop()
		return
	}

	go func() {
		for {
			select {
			case <-watcher.stop:
				return
			case <-time.After(time.Until(time.Unix(next.FiresAt, 0))):
				for _, callback := range watcher.callbacks {
					callback(next)
				}
				if err := watcher.Repo.DeleteByID(context.Background(), next.ID); err != nil {
					log.Error().Err(err).Msg("could not delete reminder")
				}
				watcher.Reset()
				return
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
	watcher.Stop()
	watcher.Start()
}

// Subscribe registers a callback for reminder fire events
func (watcher *Watcher) Subscribe(callback Callback) {
	watcher.callbacks = append(watcher.callbacks, callback)
}
