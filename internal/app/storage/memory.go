package storage

import (
	"log"
	"sync"
)

type MemoryStorage struct {
	AliasKeysMap map[string]string
	URLKeysMap   map[string]string
}

type SyncMemoryStorage struct {
	Mu            sync.Mutex
	MemoryStorage *MemoryStorage
}

func NewMemoryStorage() *SyncMemoryStorage {
	return &SyncMemoryStorage{
		Mu: sync.Mutex{},
		MemoryStorage: &MemoryStorage{
			AliasKeysMap: make(map[string]string),
			URLKeysMap:   make(map[string]string),
		},
	}
}

func (s *SyncMemoryStorage) Add(url, alias string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.MemoryStorage.AliasKeysMap[alias] = url
	s.MemoryStorage.URLKeysMap[url] = alias
	return nil
}

func (s *SyncMemoryStorage) GetURL(alias string) (url string, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if url, ok := s.MemoryStorage.AliasKeysMap[alias]; !ok {
		log.Println("URL by alias " + alias + " is not exists")
		return "", err
	} else {
		return url, nil
	}
}

func (s *SyncMemoryStorage) GetAlias(url string) (alias string, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if alias, ok := s.MemoryStorage.URLKeysMap[url]; !ok {
		log.Println("Alias by URL " + url + " is not exists")
		return "", err
	} else {
		return alias, nil
	}
}
