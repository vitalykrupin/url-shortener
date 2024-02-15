package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v10"
)

const (
	defaultServerAddress   = "localhost:8080"
	defaultResponseAddress = "http://localhost:8080"
	defaultDBDSN           = ""
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	ResponseAddress string `env:"BASE_URL"`
	FileStorePath   string `env:"FILE_STORAGE_PATH"`
	DBDSN           string `env:"DATABASE_DSN"`
}

func (conf *Config) InitConfig() {
	conf.ServerAddress = defaultServerAddress
	conf.ResponseAddress = defaultResponseAddress
	conf.DBDSN = defaultDBDSN
	conf.FileStorePath = os.TempDir() + "short-url-db.json"
	conf.parseFlags()
	conf.parseEnv()
}

func (conf *Config) parseEnv() {
	err := env.Parse(conf)
	if err != nil {
		fmt.Println("Config is not available", err)
	}
}

func (conf *Config) parseFlags() {
	flag.Func("a", "example: '-a localhost:8080'", func(addr string) error {
		conf.ServerAddress = addr
		return nil
	})
	flag.Func("b", "example: '-b http://localhost:8000'", func(addr string) error {
		conf.ResponseAddress = addr
		return nil
	})
	flag.Func("f", "example: '-f /tmp/testfile.json'", func(path string) error {
		conf.FileStorePath = path
		return nil
	})
	flag.Func("d", "example: '-d postgres://postgres:pwd@localhost:5432/postgres?sslmode=disable'", func(dbAddr string) error {
		conf.DBDSN = dbAddr
		return nil
	})
	flag.Parse()
}
