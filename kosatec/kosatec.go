package kosatec

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/computerextra/sw6-product-sync/config"
	"github.com/gocarina/gocsv"
)

type KostecProduct struct {
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

func ReadFile(path string, conf config.Config) ([]*KostecProduct, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	products := []*KostecProduct{}
	del := []rune(conf.Delimiter.Kosatec)

	// Change CSV Reader to correct Delimiter
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.LazyQuotes = true
		r.Comma = del[0]
		return r
	})

	if err := gocsv.UnmarshalFile(file, &products); err != nil {
		return nil, err
	}

	products, err = sort_products(products, conf)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func sort_products(products []*KostecProduct, conf config.Config) ([]*KostecProduct, error) {
	// TODO: Implement Sorting!

	return products, nil
}
