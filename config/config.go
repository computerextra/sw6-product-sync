package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Delimiter  Delimiter
	Autschlag  Autschlag
	Ignore     Ignore
	Override   []Override
	Kategorien []Kategorien
	UVP        []UVP
}

type Delimiter struct {
	Kosatec  rune `toml:"Delimiter.Kosatec"`
	Wortmann rune `toml:"Delimiter.Wortmann"`
}

type Autschlag struct {
	Protenzual int `toml:"Aufschlag.Protenzual"`
	Kategorie  []AufschlagKategorie
}

type AufschlagKategorie struct {
	Kategorie string `toml:"Aufschlag.Kategorie.Kategorie"`
	Prozent   int    `toml:"Aufschlag.Kategorie.Prozent"`
}

type Ignore struct {
	Kategorien []string `toml:"Ignore.Kategorien"`
	Produkte   []string `toml:"Ignore.Produkte"`
}

type Override struct {
	AlterName string `toml:"Override.AlterName"`
	NeuerName string `toml:"Override.NeuerName"`
	Index     int    `toml:"Override.Index"`
}

type Kategorien struct {
	Name            string   `toml:"Kategorien.Name"`
	Unterkategorien []string `toml:"Kategorien.Unterkategorien"`
}

type UVP struct {
	Artikelnummer string  `toml:"UVP.Artikelnummer"`
	Brutto        float32 `toml:"UVP.Brutto"`
	Netto         float32 `toml:"UVP.Netto"`
}

func (c Config) load() (*Config, error) {
	f := "config.toml"
	if _, err := os.Stat(f); err != nil {
		panic(err)
	}
	var conf Config
	_, err := toml.Decode(f, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func New() (*Config, error) {
	var conf Config
	return conf.load()
}
