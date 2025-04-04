package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Aufschlag struct {
	Kategorie string
	Prozent   int
}

type Config struct {
	Delimiter struct {
		Kosatec  string
		Wortmann string
	}
	Aufschlag struct {
		Prozentual int
	}
	Kategorie map[string][]Aufschlag
	Ignore    struct {
		Kategorien []string
		Produkte   []string
	}
	Override []struct {
		AlterName string
		NeuerName string
		Index     int
	}
	Kategorien []struct {
		Name            string
		Unterkategorien []string
	}
	UVP []struct {
		Artikelnummer string
		Brutto        float64
		Netto         float64
	}
}

func (c Config) load() (*Config, error) {
	f := "config.toml"
	if _, err := os.Stat(f); err != nil {
		panic(err)
	}
	var conf Config
	_, err := toml.DecodeFile(f, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func New() (*Config, error) {
	conf := Config{}
	return conf.load()
}
