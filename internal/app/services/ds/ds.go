package ds

import (
	"context"
	"sync"

	"github.com/vitalykrupin/url-shortener/internal/app/storage"
)

type payload struct {
	userID string
	urls   []string
}

type DeleteServiceInterface interface {
	Add(userID string, urls []string)
	Start(workers int)
	Stop()
}

type DeleteService struct {
	store storage.Storage
	input chan payload
	wg    sync.WaitGroup
}

func NewDeleteService(store storage.Storage) *DeleteService {
	return &DeleteService{
		store: store,
		input: make(chan payload),
	}
}

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

func (ds *DeleteService) Stop() {
	close(ds.input)
	ds.wg.Wait()
}

func (ds *DeleteService) Add(userID string, urls []string) {
	ds.input <- payload{
		userID: userID,
		urls:   urls,
	}
}
