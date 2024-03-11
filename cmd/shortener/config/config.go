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

func NewConfig() *Config {
	return &Config{
		ServerAddress:   defaultServerAddress,
		ResponseAddress: defaultResponseAddress,
		FileStorePath:   os.TempDir() + "short-url-db.json",
		DBDSN:           defaultDBDSN,
	}
}

func (c *Config) ParseFlags() {
	flag.Func("a", "example: '-a localhost:8080'", func(addr string) error {
		c.ServerAddress = addr
		return nil
	})
	flag.Func("b", "example: '-b http://localhost:8000'", func(addr string) error {
		c.ResponseAddress = addr
		return nil
	})
	flag.Func("f", "example: '-f /tmp/testfile.json'", func(path string) error {
		c.FileStorePath = path
		return nil
	})
	flag.Func("d", "example: '-d postgres://postgres:pwd@localhost:5432/postgres?sslmode=disable'", func(dbAddr string) error {
		c.DBDSN = dbAddr
		return nil
	})
	flag.Parse()

	err := env.Parse(c)
	if err != nil {
		fmt.Println("Config is not available", err)
	}
}
