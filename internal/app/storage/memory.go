package storage

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

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

func (storage *Storage) Add(url, alias string) error {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()
	storage.MemoryStorage.AliasKeysMap[alias] = url
	storage.MemoryStorage.URLKeysMap[url] = alias
	return nil
}

func (storage *Storage) GetURL(alias string) (url string, err error) {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()
	if url, ok := storage.MemoryStorage.AliasKeysMap[alias]; !ok {
		log.Println("URL by alias " + alias + " is not exists")
		return "", err
	} else {
		return url, nil
	}
}

func (storage *Storage) GetAlias(url string) (alias string, err error) {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()
	if alias, ok := storage.MemoryStorage.URLKeysMap[url]; !ok {
		log.Println("Alias by URL " + url + " is not exists")
		return "", err
	} else {
		return alias, nil
	}
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
