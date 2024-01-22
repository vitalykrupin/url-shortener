package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v10"
)

const (
	defaultServerAddress = "localhost:8080"
	defaultResponseAddress = "http://localhost:8080"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	ResponseAddress string `env:"BASE_URL"`
}

func (conf *Config) InitConfig() {
	conf.ServerAddress = defaultServerAddress
	conf.ResponseAddress = defaultResponseAddress
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
