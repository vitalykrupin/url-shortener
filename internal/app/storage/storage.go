package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type StorageKeeper interface {
	AddToMemoryStore(url, alias string)
	GetURL(alias string) (url string, ok bool)
	GetAlias(url string) (alias string, ok bool)
	SaveJSONtoFS(path string)
	LoadJSONfromFS(path string) error
}

type MemoryStorage struct {
	AliasKeysMap map[string]string
	URLKeysMap   map[string]string
}

type Storage struct {
	Mu            sync.Mutex
	MemoryStorage MemoryStorage
}

func NewMemoryStorage() *Storage {
	storage := new(Storage)
	storage.MemoryStorage.AliasKeysMap = make(map[string]string)
	storage.MemoryStorage.URLKeysMap = make(map[string]string)
	storage.Mu = sync.Mutex{}
	return storage
}

func (storage *Storage) AddToMemoryStore(url, alias string) {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()
	storage.MemoryStorage.AliasKeysMap[alias] = url
	storage.MemoryStorage.URLKeysMap[url] = alias
}

func (storage *Storage) GetURL(alias string) (url string, found bool) {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()
	url, found = storage.MemoryStorage.AliasKeysMap[alias]
	return url, found
}

func (storage *Storage) GetAlias(url string) (alias string, found bool) {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()
	alias, found = storage.MemoryStorage.URLKeysMap[url]
	return alias, found
}

func (storage *Storage) SaveJSONtoFS(path string) {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()
	if path == "" {
		return
	}
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	jsonData, err := json.MarshalIndent(storage.MemoryStorage.AliasKeysMap, "", "	")
	if err != nil {
		panic(err)
	}
	_, err = file.Write(jsonData)
	if err != nil {
		panic(err)
	}
}

func (storage *Storage) LoadJSONfromFS(path string) error {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()
	if path == "" {
		return nil
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &storage.MemoryStorage.AliasKeysMap)
	if err != nil {
		return err
	}
	for k, v := range storage.MemoryStorage.AliasKeysMap {
		storage.MemoryStorage.URLKeysMap[v] = k
	}

	return nil
}
