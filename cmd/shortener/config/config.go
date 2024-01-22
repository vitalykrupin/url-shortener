package config

import (
	"flag"
)

type Config struct {
	ServerAddress   string
	ResponseAddress string
}

func (conf *Config) ParseFlags() {
	conf.ServerAddress = "localhost:8080"
	conf.ResponseAddress = "http://localhost:8080"
	flag.Func("a", "example '-a localhost:8080'", func(addr string) error {
		conf.ServerAddress = addr
		return nil
	})
	flag.Func("b", "example '-b http://localhost:8000'", func(addr string) error {
		conf.ResponseAddress = addr
		return nil
	})
	flag.Parse()
}
