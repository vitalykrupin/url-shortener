package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type JSONFS struct {
	UUID  string      `json:"id"`
	Alias Alias       `json:"alias"`
	URL   OriginalURL `json:"url"`
}

type FileStorage struct {
	SyncMemoryStorage *SyncMemoryStorage
	file              *os.File
}

func NewFileStorage(FileStoragePath string) (*FileStorage, error) {
	if FileStoragePath == "" {
		return nil, fmt.Errorf("no FileStoragePath provided")
	}
	syncMem := NewMemoryStorage()

	file, err := os.OpenFile(FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("can not open file: %w", err)
	}
	fs := FileStorage{SyncMemoryStorage: syncMem, file: file}
	if err := fs.LoadJSONfromFS(); err != nil && !errors.Is(err, bufio.ErrTooLong) {
		return nil, fmt.Errorf("can not load JSON from file: %w", err)
	}
	return &fs, nil
}

func (f *FileStorage) LoadJSONfromFS() error {
	if _, err := f.file.Seek(0, 0); err != nil {
		return err
	}
	scanner := bufio.NewScanner(f.file)
	data := make(map[Alias]OriginalURL)
	for scanner.Scan() {
		var urls JSONFS
		if err := json.Unmarshal(scanner.Bytes(), &urls); err != nil {
			return err
		}
		data[urls.Alias] = urls.URL
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := f.SyncMemoryStorage.Add(data); err != nil {
		return err
	}
	return nil
}

func (f *FileStorage) Add(ctx context.Context, batch map[Alias]OriginalURL) error {
	if err := f.SyncMemoryStorage.Add(batch); err != nil {
		return err
	}
	if f.file == nil {
		return errors.New("file is not opened")
	}

	writter := bufio.NewWriter(f.file)
	for alias, url := range batch {
		entry := JSONFS{
			UUID:  strconv.Itoa(len(f.SyncMemoryStorage.MemoryStorage.AliasKeysMap)),
			Alias: alias,
			URL:   url,
		}
		data, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		if _, err := writter.Write(data); err != nil {
			return err
		}
		if err := writter.WriteByte('\n'); err != nil {
			return err
		}
	}
	if err := writter.Flush(); err != nil {
		return err
	}
	return nil
}

func (f *FileStorage) GetURL(ctx context.Context, alias Alias) (url OriginalURL, err error) {
	return f.SyncMemoryStorage.GetURL(alias)
}

func (f *FileStorage) GetUserURLs(ctx context.Context, userID string) (aliasKeysMap AliasKeysMap, err error) {
	return nil, fmt.Errorf("can not get user urls from file storage")
}

func (f *FileStorage) GetAlias(ctx context.Context, url OriginalURL) (alias Alias, err error) {
	return f.SyncMemoryStorage.GetAlias(url)
}

func (f *FileStorage) DeleteUserURLs(ctx context.Context, userID string, urls []string) error {
	return fmt.Errorf("can not delete user urls from file storage")
}

func (f *FileStorage) CloseStorage(ctx context.Context) error {
	if f.file != nil {
		return f.file.Close()
	}
	return nil
}

func (f *FileStorage) PingStorage(ctx context.Context) error { return nil }
