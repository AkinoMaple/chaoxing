package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Account struct {
		Cookie string `toml:"cookie"`
	} `toml:"account"`
	Settings struct {
		Routines int     `toml:"routines"`
		FetchID  [][]int `toml:"fetch_id"`
	} `toml:"settings"`
	MongoDB struct {
		ApplyURI   string `toml:"apply_uri"`
		Database   string `toml:"database"`
		Collection string `toml:"collection"`
	} `toml:"mongodb"`
}

func (c *Config) Check() {
	if c.Account.Cookie == "" {
		log.Panicln("Cookie can not be null.")
	}

	if c.Settings.Routines < 1 {
		log.Panicln("Routines can not less than 1.")
	}

	if c.Settings.FetchID == nil {
		log.Panicln("Pleace set Fetch ID,")
	}
}

func unmarshal(c string) *Config {
	var config Config
	if _, err := toml.Decode(c, &config); err != nil {
		log.Panic(err)
	}
	config.Check()

	return &config
}

func loadConfig(cfp *string) string {
	f, err := os.Open(*cfp)
	if err != nil {
		log.Panic(err)
	}

	defer f.Close()

	buf := make([]byte, 1024)
	var config string

	for n, err := f.Read(buf); err == nil; n, err = f.Read(buf) {
		config += string(buf[:n])
	}

	return config
}
