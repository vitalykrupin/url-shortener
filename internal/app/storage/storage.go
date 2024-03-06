package storage

import (
	"context"

	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
)

type Alias string
type OriginalURL string

type Storage interface {
	Add(ctx context.Context, batch map[Alias]OriginalURL) error
	GetURL(ctx context.Context, alias Alias) (url OriginalURL, err error)
	GetAlias(ctx context.Context, url OriginalURL) (alias Alias, err error)
	GetUserURLs(ctx context.Context, userID string) (aliasKeysMap AliasKeysMap, err error)
	DeleteUserURLs(ctx context.Context, userID string, urls []string) error
	CloseStorage(ctx context.Context) error
	PingStorage(ctx context.Context) error
}

var Store Storage

func NewStorage(conf *config.Config) (Storage, error) {
	if conf.DBDSN != "" {
		return NewDB(conf.DBDSN)
	} else {
		return NewFileStorage(conf.FileStorePath)
	}
}
