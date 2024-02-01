package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	AliasKeysMap   map[string]string
	FullURLKeysMap map[string]string
}

type Store struct {
	Mu    sync.Mutex
	Store DB
}

func NewStorage() *Store {
	store := new(Store)
	store.Store.AliasKeysMap = make(map[string]string)
	store.Store.FullURLKeysMap = make(map[string]string)
	store.Mu = sync.Mutex{}
	return store
}

func (store *Store) SaveJSONtoFS(path string) {
	store.Mu.Lock()
	defer store.Mu.Unlock()
	if path == "" {
		return
	}
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	jsonData, err := json.MarshalIndent(store.Store.AliasKeysMap, "", "	")
	if err != nil {
		panic(err)
	}
	_, err = file.Write(jsonData)
	if err != nil {
		panic(err)
	}
}

func (store *Store) LoadJSONfromFS(path string) error {
	store.Mu.Lock()
	defer store.Mu.Unlock()
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
	err = json.Unmarshal(data, &store.Store.AliasKeysMap)
	if err != nil {
		return err
	}
	for k, v := range store.Store.AliasKeysMap {
		store.Store.FullURLKeysMap[v] = k
	}

	return nil
}
