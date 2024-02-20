package storage

import "context"

type Alias string
type OriginalURL string

type Storage interface {
	Add(ctx context.Context, batch map[Alias]OriginalURL) error
	GetURL(ctx context.Context, alias Alias) (url OriginalURL, err error)
	GetAlias(ctx context.Context, url OriginalURL) (alias Alias, err error)
	CloseStorage(ctx context.Context) error
	PingStorage(ctx context.Context) error
}
