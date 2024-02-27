package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
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

func NewFileStorage(cfg *config.Config) Storage {
	if cfg.FileStorePath == "" {
		return nil
	}
	syncMem := NewMemoryStorage()

	file, err := os.Create(cfg.FileStorePath)
	if err != nil {
		log.Fatal("Can not create file")
	}
	fs := &FileStorage{SyncMemoryStorage: syncMem, file: file}
	err = fs.LoadJSONfromFS()
	if err != nil {
		log.Fatal("Can not load JSON from file")
	}
	return fs
}

func (f *FileStorage) LoadJSONfromFS() error {
	scanner := bufio.NewScanner(f.file)
	data := make(map[Alias]OriginalURL)
	for scanner.Scan() {
		urls := JSONFS{}
		err := json.Unmarshal(scanner.Bytes(), &urls)
		if err != nil {
			return err
		}
		data[urls.Alias] = urls.URL
	}
	err := f.SyncMemoryStorage.Add(data)
	if err != nil {
		log.Println("Can not save data to memory")
		return err
	}
	return nil
}

func (f *FileStorage) Add(ctx context.Context, batch map[Alias]OriginalURL) error {
	err := f.SyncMemoryStorage.Add(batch)
	if err != nil {
		log.Println("Can not save data to memory")
		return err
	}
	if f.file == nil {
		log.Println("Can not save data to file, file is not exists")
		return nil
	}

	writter := bufio.NewWriter(f.file)
	var urls = []JSONFS{}
	for alias, url := range batch {
		urls = append(urls, JSONFS{
			UUID:  strconv.Itoa(len(f.SyncMemoryStorage.MemoryStorage.AliasKeysMap)),
			Alias: alias,
			URL:   url,
		})
	}
	data, err := json.Marshal(urls)
	if err != nil {
		log.Println("Can not marshal data")
		return nil
	}
	if _, err := writter.Write(data); err != nil {
		log.Println("Can not write data into file")
		return nil
	}

	if err := writter.WriteByte('\n'); err != nil {
		log.Println("Can not write new line into file")
		return nil
	}

	if err := writter.Flush(); err != nil {
		log.Println("Can not flush data into file")
		return nil
	}

	return nil
}

func (f *FileStorage) GetURL(ctx context.Context, alias Alias) (url OriginalURL, err error) {
	return f.SyncMemoryStorage.GetURL(alias)
}

func (f *FileStorage) GetUserURLs(ctx context.Context, userID string) (aliasKeysMap *aliasKeysMap, err error) {
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
