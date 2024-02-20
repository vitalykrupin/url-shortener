package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
)

type JSONFS struct {
	UUID  string `json:"id"`
	Alias string `json:"alias"`
	URL   string `json:"url"`
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
	data := make(map[string]string)
	for scanner.Scan() {
		urls := JSONFS{}
		err := json.Unmarshal(scanner.Bytes(), &urls)
		if err != nil {
			return err
		}
		data[urls.Alias] = urls.URL
	}
	for alias, url := range data {
		err := f.SyncMemoryStorage.Add(alias, url)
		if err != nil {
			log.Println("Can not save data to memory")
			return err
		}
	}
	return nil
}

func (f *FileStorage) Add(ctx context.Context, alias, url string) error {
	err := f.SyncMemoryStorage.Add(alias, url)
	if err != nil {
		log.Println("Can not save data to memory")
		return err
	}
	if f.file == nil {
		log.Println("Can not save data to file, file is not exists")
		return nil
	}

	writter := bufio.NewWriter(f.file)
	urls := JSONFS{
		UUID:  strconv.Itoa(len(f.SyncMemoryStorage.MemoryStorage.AliasKeysMap)),
		Alias: alias,
		URL:   url,
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

func (f *FileStorage) GetURL(ctx context.Context, alias string) (url string, err error) {
	return f.SyncMemoryStorage.GetURL(alias)
}

func (f *FileStorage) GetAlias(ctx context.Context, url string) (alias string, err error) {
	return f.SyncMemoryStorage.GetAlias(url)
}

func (f *FileStorage) CloseStorage(ctx context.Context) error {
	if f.file != nil {
		return f.file.Close()
	}
	return nil
}

func (f *FileStorage) PingStorage(ctx context.Context) error { return nil }
