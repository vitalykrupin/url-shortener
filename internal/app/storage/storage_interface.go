package storage

type StorageInterface interface {
	Add(url, alias string) error
	GetURL(alias string) (url string, err error)
	GetAlias(url string) (alias string, err error)
	SaveJSONtoFS(path string)
	LoadJSONfromFS(path string) error
}
