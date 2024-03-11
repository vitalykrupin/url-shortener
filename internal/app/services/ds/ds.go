package ds

import (
	"context"

	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

type payload struct {
	userID string
	urls   []string
}

type DeleteService struct {
	store storage.Storage
	input chan payload
}

func NewDeleteService(store storage.Storage) *DeleteService {
	return &DeleteService{
		store: store,
		input: make(chan payload),
	}
}

func (ds *DeleteService) Start(workers int) {
	result := make(chan payload)
	for w := 1; w <= workers; w++ {
		go func() {
			for p := range ds.input {
				result <- p
			}
		}()
	}
	go func() {
		for p := range result {
			ds.store.DeleteUserURLs(context.Background(), p.userID, p.urls)
		}
	}()
}

func (ds *DeleteService) Stop() {
	close(ds.input)
}

func (ds *DeleteService) Add(userID string, urls []string) {
	ds.input <- payload{
		userID: userID,
		urls:   urls,
	}
}
