package kosatec

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/computerextra/sw6-product-sync/config"
	"github.com/computerextra/sw6-product-sync/shopware"
	"github.com/gocarina/gocsv"
)

type KosatecProduct struct {
	Artnr          string `csv:"artnr"`
	Herstnr        string `csv:"herstnr"`
	Artname        string `csv:"artname"`
	Hersteller     string `csv:"hersteller"`
	Hersturl       string `csv:"hersturl"`
	Ean            string `csv:"ean"`
	Hek            string `csv:"hek"` // decimal point: "."
	Vkbrutto       string `csv:"vkbrutto"`
	Verfuegbar     string `csv:"verfuegbar"`
	Menge          string `csv:"menge"`
	Eta            string `csv:"eta"`
	Indate         string `csv:"indate"`
	Gewicht        string `csv:"gewicht"`
	Eol            string `csv:"eol"`
	Kat1           string `csv:"kat1"`
	Kat2           string `csv:"kat2"`
	Kat3           string `csv:"kat3"`
	Kat4           string `csv:"kat4"`
	Kat5           string `csv:"kat5"`
	Kat6           string `csv:"kat6"`
	Title          string `csv:"title"`
	Short_desc     string `csv:"short_desc"`
	Short_summary  string `csv:"short_summary"`
	Long_summary   string `csv:"long_summary"`
	Marketing_text string `csv:"marketing_text"`
	Specs          string `csv:"specs"`
	Pdf            string `csv:"pdf"`
	Pdf_manual     string `csv:"pdf_manual"`
	Images_s       string `csv:"images_s"`  // delimiter ";"
	Images_m       string `csv:"images_m"`  // delimiter ";"
	Images_l       string `csv:"images_l"`  // delimiter ";"
	Images_xl      string `csv:"images_xl"` // delimiter ";"
}

func ReadFile(path string, conf config.Config) ([]shopware.Artikel, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	products := []KosatecProduct{}

	// TODO: Weg finden, um den Delimiter aus der Toml Datei zu nutzen

	// Change CSV Reader to correct Delimiter
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.LazyQuotes = true
		r.Comma = '\t'
		return r
	})

	if err := gocsv.UnmarshalFile(file, &products); err != nil {
		return nil, err
	}

	return sort_products(products, conf)
}

func sort_products(products []KosatecProduct, conf config.Config) ([]shopware.Artikel, error) {
	var sorted []shopware.Artikel

	for _, item := range products {
		stop := false
		if !skip(item, conf) {
			a := shopware.Artikel{}
			a.Artikelnummer = fmt.Sprintf("K%s", strings.TrimSpace(item.Artnr))
			if len(item.Menge) > 1 {
				bestand, err := strconv.Atoi(item.Menge)
				if err != nil {
					bestand = 0
				}
				a.Bestand = bestand
			} else {
				a.Bestand = 0
			}
			a.HerstellerNummer = strings.TrimSpace(item.Herstnr)
			a.Name = strings.TrimSpace(item.Artname)
			a.Ean = strings.TrimSpace(item.Ean)
			a.Beschreibung = fmt.Sprintf("%s<br>%s", item.Marketing_text, item.Long_summary)
			if len(item.Hersteller) > 1 {
				a.Hersteller = strings.TrimSpace(item.Hersteller)
			} else {
				stop = true
			}
			if len(item.Kat1) > 1 {
				a.Kategorie1 = check_category(item.Kat1, 1, conf)
			} else {
				stop = true
			}
			if len(item.Kat2) > 1 {
				a.Kategorie2 = check_category(item.Kat2, 2, conf)
			}
			if len(item.Kat3) > 1 {
				a.Kategorie3 = check_category(item.Kat3, 3, conf)
			}
			if len(item.Kat4) > 1 {
				a.Kategorie4 = check_category(item.Kat4, 4, conf)
			}
			if len(item.Kat5) > 1 {
				a.Kategorie5 = check_category(item.Kat5, 5, conf)
			}
			if len(item.Kat6) > 1 {
				a.Kategorie6 = check_category(item.Kat6, 6, conf)
			}
			ekFloat, err := strconv.ParseFloat(item.Hek, 64)
			if err != nil {
				ekFloat = 0
				stop = true
			}
			a.Ek = ekFloat
			if len(item.Images_xl) > 1 {
				a.Bilder = strings.ReplaceAll(item.Images_xl, ";", "|")
			} else if len(item.Images_l) > 1 {
				a.Bilder = strings.ReplaceAll(item.Images_l, ";", "|")
			} else if len(item.Images_m) > 1 {
				a.Bilder = strings.ReplaceAll(item.Images_m, ";", "|")
			} else if len(item.Images_s) > 1 {
				a.Bilder = strings.ReplaceAll(item.Images_s, ";", "|")
			}
			if !stop {
				sorted = append(sorted, a)
			}
		}
	}

	return sorted, nil
}

func skip(item KosatecProduct, conf config.Config) bool {
	if len(item.Artnr) < 1 {
		return true
	}
	if len(item.Kat1) < 1 {
		return true
	}
	if len(item.Hek) < 1 {
		return true
	}
	if len(item.Hersteller) < 1 {
		return true
	}
	if len(item.Artname) < 1 {
		return true
	}
	if is_ignored(conf.Ignore.Produkte, item.Artnr) {
		return true
	}
	if is_ignored(conf.Ignore.Kategorien, item.Kat1) {
		return true
	}
	if is_ignored(conf.Ignore.Kategorien, item.Kat2) {
		return true
	}
	if is_ignored(conf.Ignore.Kategorien, item.Kat3) {
		return true
	}
	if is_ignored(conf.Ignore.Kategorien, item.Kat4) {
		return true
	}
	if is_ignored(conf.Ignore.Kategorien, item.Kat5) {
		return true
	}
	if is_ignored(conf.Ignore.Kategorien, item.Kat6) {
		return true
	}
	return false
}

func is_ignored(ignored []string, str string) bool {
	if len(ignored) < 1 {
		return false
	}
	if len(str) < 1 {
		return false
	}

	for _, x := range ignored {
		if x == str {
			return true
		}
		if strings.TrimSpace(str) == x {
			return true
		}
		if strings.TrimSpace(x) == str {
			return true
		}
		if strings.TrimSpace(x) == strings.TrimSpace(str) {
			return true
		}
	}
	return false
}

func check_category(cat string, idx int, conf config.Config) string {
	strippped := strings.TrimSpace(cat)
	if len(conf.Override) < 1 {
		return strippped
	}
	for _, x := range conf.Override {
		if strippped == x.AlterName && idx == x.Index {
			return x.NeuerName
		}
	}
	return strippped
}
