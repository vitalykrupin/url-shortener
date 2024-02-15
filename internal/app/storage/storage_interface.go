package storage

import "context"

type Storage interface {
	Add(ctx context.Context, url, alias string) error
	GetURL(ctx context.Context, alias string) (url string, err error)
	GetAlias(ctx context.Context, url string) (alias string, err error)
	CloseStorage(ctx context.Context) error
	PingStorage(ctx context.Context) error
	// SaveJSONtoFS(path string)
	// LoadJSONfromFS(path string) error
}
