package wortmann

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

type ProductCatalog struct {
	ProductId                        string `csv:"ProductId,omitempty"`
	ReferenceNo                      string `csv:"ReferenceNo,omitempty"`
	EAN                              string `csv:"EAN,omitempty"`
	Manufacturer                     string `csv:"Manufacturer,omitempty"`
	Price_B2B_Regular                string `csv:"Price_B2B_Regular,omitempty"`
	Price_B2B_Discounted             string `csv:"Price_B2B_Discounted,omitempty"`
	Price_B2B_DiscountPercent        string `csv:"Price_B2B_DiscountPercent,omitempty"`
	Price_B2B_DiscountAmount         string `csv:"Price_B2B_DiscountAmount,omitempty"`
	Price_B2C_exclVAT                string `csv:"Price_B2C_exclVAT,omitempty"`
	Price_B2C_inclVAT                string `csv:"Price_B2C_inclVAT,omitempty"`
	Price_B2C_VATRate                string `csv:"Price_B2C_VATRate,omitempty"`
	Price_B2C_VATCountry             string `csv:"Price_B2C_VATCountry,omitempty"`
	Price_B2X_Currency               string `csv:"Price_B2X_Currency,omitempty"`
	Stock                            int64  `csv:"Stock,omitempty"`
	StockNextDelivery                string `csv:"StockNextDelivery,omitempty"`
	StockNextDeliveryAccessVolume    string `csv:"StockNextDeliveryAccessVolume,omitempty"`
	WarrantyCode                     string `csv:"WarrantyCode,omitempty"`
	EOL                              string `csv:"EOL,omitempty"`
	Promotion                        string `csv:"Promotion,omitempty"`
	NonReturnable                    string `csv:"NonReturnable,omitempty"`
	RemainingStock                   string `csv:"RemainingStock,omitempty"`
	ImagePrimary                     string `csv:"ImagePrimary,omitempty"`
	ImageAdditional                  string `csv:"ImageAdditional,omitempty"`
	ProductLink                      string `csv:"ProductLink,omitempty"`
	GrossWeight                      string `csv:"GrossWeight,omitempty"`
	NetWeight                        string `csv:"NetWeight,omitempty"`
	RelatedProducts                  string `csv:"RelatedProducts,omitempty"`
	AccessoryProducts                string `csv:"AccessoryProducts,omitempty"`
	Description_1031_German          string `csv:"Description_1031_German,omitempty"`
	CategoryName_1031_German         string `csv:"CategoryName_1031_German,omitempty"`
	CategoryPath_1031_German         string `csv:"CategoryPath_1031_German,omitempty"`
	WarrantyDescription_1031_German  string `csv:"WarrantyDescription_1031_German,omitempty"`
	Description_1033_English         string `csv:"Description_1033_English,omitempty"`
	CategoryName_1033_English        string `csv:"CategoryName_1033_English,omitempty"`
	CategoryPath_1033_English        string `csv:"CategoryPath_1033_English,omitempty"`
	WarrantyDescription_1033_English string `csv:"WarrantyDescription_1033_English,omitempty"`
	Description_1036_French          string `csv:"Description_1036_French,omitempty"`
	CategoryName_1036_French         string `csv:"CategoryName_1036_French,omitempty"`
	CategoryPath_1036_French         string `csv:"CategoryPath_1036_French,omitempty"`
	WarrantyDescription_1036_French  string `csv:"WarrantyDescription_1036_French,omitempty"`
	Description_1043_Dutch           string `csv:"Description_1043_Dutch,omitempty"`
	CategoryName_1043_Dutch          string `csv:"CategoryName_1043_Dutch,omitempty"`
	CategoryPath_1043_Dutch          string `csv:"CategoryPath_1043_Dutch,omitempty"`
	WarrantyDescription_1043_Dutch   string `csv:"WarrantyDescription_1043_Dutch,omitempty"`
	ProductDisplayType               string `csv:"ProductDisplayType,omitempty"`
	LicenseTypeCode                  string `csv:"LicenseTypeCode,omitempty"`
	LicenseTypeDescription           string `csv:"LicenseTypeDescription,omitempty"`
}

type Content struct {
	ProductId                    string `csv:"ProductId,omitempty"`
	PrintText_1031_German        string `csv:"PrintText_1031_German,omitempty"`
	PrintText_1033_English       string `csv:"PrintText_1033_English,omitempty"`
	PrintText_1036_French        string `csv:"PrintText_1036_French,omitempty"`
	PrintText_1043_Dutch         string `csv:"PrintText_1043_Dutch,omitempty"`
	LongDescription_1031_German  string `csv:"LongDescription_1031_German,omitempty"`
	LongDescription_1033_English string `csv:"LongDescription_1033_English,omitempty"`
	LongDescription_1036_French  string `csv:"LongDescription_1036_French,omitempty"`
	LongDescription_1043_Dutch   string `csv:"LongDescription_1043_Dutch,omitempty"`
}

func ReadFile(path1 string, path2 string, conf config.Config) ([]shopware.Artikel, error) {
	f1, err := os.OpenFile(path1, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f1.Close()
	f2, err := os.OpenFile(path2, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f2.Close()

	catalog := []ProductCatalog{}
	content := []Content{}

	// TODO: Weg finden, um den Delimiter aus der Toml Datei zu nutzen

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.LazyQuotes = true
		r.Comma = ';'
		return r
	})
	if err := gocsv.UnmarshalFile(f1, &catalog); err != nil {
		return nil, err
	}
	if err := gocsv.UnmarshalFile(f2, &content); err != nil {
		return nil, err
	}

	return sort_products(catalog, content, conf)
}

func sort_products(catalog []ProductCatalog, content []Content, conf config.Config) ([]shopware.Artikel, error) {
	var sorted []shopware.Artikel

	for _, item := range catalog {
		stop := false
		if !skip(item, conf) {
			a := shopware.Artikel{}
			a.Artikelnummer = fmt.Sprintf("W%s", strings.TrimSpace(item.ProductId))
			a.Hersteller = "WORTMANN AG"
			a.Name = item.Description_1031_German
			a.Ean = strings.TrimSpace(item.EAN)
			priceF, err := strconv.ParseFloat(strings.TrimSpace(item.Price_B2C_inclVAT), 64)
			if err != nil {
				priceF = 0
				stop = true
			}
			a.Vk = priceF
			a.Bestand = item.Stock
			bilder := strings.TrimSpace(item.ImagePrimary)
			if len(strings.TrimSpace(item.ImageAdditional)) > 1 {
				bilder = fmt.Sprintf("%s|%s", bilder, strings.TrimSpace(item.ImageAdditional))
			}
			a.Bilder = bilder
			kat := strings.TrimSpace(item.CategoryName_1031_German)
			a.Kategorie1 = check_cat(kat)
			a = add_content(a, content)
			if !stop {
				sorted = append(sorted, a)
			}
		}
	}
	return sorted, nil
}

func add_content(item shopware.Artikel, contents []Content) shopware.Artikel {
	for _, content := range contents {
		if fmt.Sprintf("W%s", strings.TrimSpace(content.ProductId)) == item.Artikelnummer {
			item.Beschreibung = content.LongDescription_1031_German
		}
	}
	return item
}

func check_cat(kat string) string {
	switch kat {
	case "PC":
		return "Marken PCs"
	case "LCD":
		return "Monitore"
	case "Dockingstations":
		return "Zubehör Notebooks"
	case "PC- & NetzwerkkamerasC":
		return "WebCams"
	case "PAD":
		return "Tablets"
	case "Taschen":
		return "Notebooktaschen"
	case "MOBILE":
		return "Notebooks"
	case "FIREWALL":
		return "Firewall"
	case "Headset & Mikro":
		return "Kopfhörer & Headsets"
	case "THINCLIENT":
		return "Mini-PC / Barebones"
	case "ALL-IN-ONE":
		return "All in One PC-Systeme"
	}

	return kat
}

func skip(item ProductCatalog, conf config.Config) bool {
	if strings.TrimSpace(item.Manufacturer) != "WORTMANN AG" {
		return true
	}
	if is_ignored(conf.Ignore.Kategorien, item.CategoryName_1031_German) {
		return true
	}
	if is_ignored(conf.Ignore.Produkte, item.ProductId) {
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
