package storage

import "context"

type Storage interface {
	Add(ctx context.Context, batch map[Alias]OriginalURL) error
	GetURL(ctx context.Context, alias Alias) (url OriginalURL, err error)
	GetAlias(ctx context.Context, url OriginalURL) (alias Alias, err error)
	GetUserURLs(ctx context.Context, userID string) (aliasKeysMap *aliasKeysMap, err error)
	CloseStorage(ctx context.Context) error
	PingStorage(ctx context.Context) error
}
