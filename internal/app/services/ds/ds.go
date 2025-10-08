// Package ds provides delete service functionality
package ds

import (
	"context"
	"sync"

	"github.com/vitalykrupin/url-shortener/internal/app/storage"
)

// payload represents the data structure for delete operations
type payload struct {
	userID string
	urls   []string
}

// DeleteServiceInterface defines the interface for delete service
type DeleteServiceInterface interface {
	// Add adds URLs to the delete queue
	Add(userID string, urls []string)
	
	// Start starts the delete service with specified number of workers
	Start(workers int)
	
	// Stop stops the delete service
	Stop()
}

// DeleteService implements the delete service functionality
type DeleteService struct {
	store storage.Storage
	input chan payload
	wg    sync.WaitGroup
}

// NewDeleteService creates a new delete service instance
func NewDeleteService(store storage.Storage) *DeleteService {
	return &DeleteService{
		store: store,
		input: make(chan payload),
	}
}

// Start starts the delete service with specified number of workers
func (ds *DeleteService) Start(workers int) {
	for w := 1; w <= workers; w++ {
		ds.wg.Add(1)
		go func() {
			defer ds.wg.Done()
			for p := range ds.input {
				_ = ds.store.DeleteUserURLs(context.Background(), p.userID, p.urls)
			}
		}()
	}
}

// Stop stops the delete service
func (ds *DeleteService) Stop() {
	close(ds.input)
	ds.wg.Wait()
}

// Add adds URLs to the delete queue
func (ds *DeleteService) Add(userID string, urls []string) {
	ds.input <- payload{
		userID: userID,
		urls:   urls,
	}
}
