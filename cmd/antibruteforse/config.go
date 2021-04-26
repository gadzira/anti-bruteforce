package main

import (
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger   LoggerConf
	Server   server
	Storage  metapart
	DataBase database
}

type LoggerConf struct {
	Level      string
	LogFile    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	LocalTime  bool
	Compress   bool
}

type server struct {
	Port string
}

type metapart struct {
	N   int
	M   int
	K   int
	TTL string
}

type database struct {
	DSN string
}

func NewConfig(fileName string) Config {
	var confDir = "../../configs" // nolint:gofumpt
	conFile := filepath.Join(confDir, fileName)

	var config Config
	if _, err := toml.DecodeFile(conFile, &config); err != nil {
		log.Fatal("Can't load configuration file:", err)
	}

	return Config{
		Server:   config.Server,
		Logger:   config.Logger,
		Storage:  config.Storage,
		DataBase: config.DataBase,
	}
}
