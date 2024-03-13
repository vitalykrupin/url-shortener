package storage

import (
	"fmt"
	"log"
	"sync"
)

type AliasKeysMap map[Alias]OriginalURL
type urlKeysMap map[OriginalURL]Alias

type MemoryStorage struct {
	AliasKeysMap AliasKeysMap
	URLKeysMap   urlKeysMap
}

type SyncMemoryStorage struct {
	Mu            sync.Mutex
	MemoryStorage *MemoryStorage
}

func NewMemoryStorage() *SyncMemoryStorage {
	return &SyncMemoryStorage{
		Mu: sync.Mutex{},
		MemoryStorage: &MemoryStorage{
			AliasKeysMap: make(map[Alias]OriginalURL),
			URLKeysMap:   make(map[OriginalURL]Alias),
		},
	}
}

func (s *SyncMemoryStorage) Add(batch map[Alias]OriginalURL) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for k, v := range batch {
		s.MemoryStorage.AliasKeysMap[k] = v
		s.MemoryStorage.URLKeysMap[v] = k
	}
	return nil
}

func (s *SyncMemoryStorage) GetURL(alias Alias) (url OriginalURL, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if url, ok := s.MemoryStorage.AliasKeysMap[alias]; !ok {
		log.Println("URL by alias " + alias + " is not exists")
		return "", fmt.Errorf("url by alias %s is not exists", alias)
	} else {
		return url, nil
	}
}

func (s *SyncMemoryStorage) GetAlias(url OriginalURL) (alias Alias, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if alias, ok := s.MemoryStorage.URLKeysMap[url]; !ok {
		log.Println("Alias by URL " + url + " is not exists")
		return "", fmt.Errorf("alias by URL %s is not exists", url)
	} else {
		return alias, nil
	}
}
