package storage

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
