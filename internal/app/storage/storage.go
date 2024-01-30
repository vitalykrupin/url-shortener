package storage

import (
	"encoding/json"
	"errors"
	"os"
)

type DB struct {
	AliasKeysMap   map[string]string
	FullURLKeysMap map[string]string
}

func NewStorage() *DB {
	db := new(DB)
	db.AliasKeysMap = make(map[string]string)
	db.FullURLKeysMap = make(map[string]string)
	return db
}

func (db *DB) SaveJSONtoFS(path string) {
	if path == "" {
		return
	}
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	jsonData, err := json.MarshalIndent(db.AliasKeysMap, "", "	")
	if err != nil {
		panic(err)
	}
	_, err = file.Write(jsonData)
	if err != nil {
		panic(err)
	}
}

func (db *DB) LoadJSONfromFS(path string) error {
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
	err = json.Unmarshal(data, &db.AliasKeysMap)
	if err != nil {
		return err
	}
	for k, v := range db.AliasKeysMap {
		db.FullURLKeysMap[v] = k
	}

	return nil
}
