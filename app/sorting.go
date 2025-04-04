package app

import (
	"log/slog"
	"strings"
	"sync"

	"github.com/computerextra/sw6-product-sync/kosatec"
	"github.com/computerextra/sw6-product-sync/shopware"
	"github.com/computerextra/sw6-product-sync/wortmann"
	sdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
)

func (a App) readKosatec() ([]shopware.Artikel, error) {
	res, err := kosatec.ReadFile(KosatecFile, *a.config)
	if err == nil {
		a.logger.Info("Successfully read Kosatec file", slog.Any("items", len(res)))
	}
	return res, err
}

func (a App) readWortmann() ([]shopware.Artikel, error) {
	res, err := wortmann.ReadFile(WortmannCatalog, WortmannContent, *a.config)
	if err == nil {
		a.logger.Info("Successfully read Wortmann file", slog.Any("items", len(res)))
	}
	return res, err
}

var Kosatec, Wortmann []shopware.Artikel
var ShopArtikel *sdk.ProductCollection
var err error

func worker(wg *sync.WaitGroup, which string, a App) {
	defer wg.Done()

	switch which {
	case "Kosatec":
		Kosatec, err = a.readKosatec()
	case "Wortmann":
		Wortmann, err = a.readWortmann()
	case "Shop":
		ShopArtikel, err = a.getAllProducts()
	}
	if err != nil {
		a.logger.Error("failed to sort products", slog.Any("error", err))
	}
}

func (a App) SortProducts() (NeueArtikel, AlteArtikel []shopware.Artikel, EolArtikel, Hersteller []string, err error) {
	var wg sync.WaitGroup

	wg.Add(1)
	go worker(&wg, "Kosatec", a)
	wg.Add(1)
	go worker(&wg, "Wortmann", a)
	wg.Add(1)
	go worker(&wg, "Shop", a)

	wg.Wait()

	for _, item := range Kosatec {
		found := false
		Hersteller = append(Hersteller, item.Hersteller)
		for _, prod := range ShopArtikel.Data {
			if item.Artikelnummer == prod.ProductNumber {
				AlteArtikel = append(AlteArtikel, item)
				found = true
				break
			}
		}
		if !found {
			NeueArtikel = append(NeueArtikel, item)
		}
	}
	for _, item := range Wortmann {
		found := false
		Hersteller = append(Hersteller, item.Hersteller)
		for _, prod := range ShopArtikel.Data {
			if item.Artikelnummer == prod.ProductNumber {
				AlteArtikel = append(AlteArtikel, item)
				found = true
				break
			}
		}
		if !found {
			NeueArtikel = append(NeueArtikel, item)
		}
	}

	for _, item := range ShopArtikel.Data {
		found := false
		for _, x := range Kosatec {
			if x.Artikelnummer == item.ProductNumber {
				found = true
				break
			}
		}
		for _, x := range Wortmann {
			if x.Artikelnummer == item.ProductNumber {
				found = true
				break
			}
		}
		if !found {
			if strings.HasPrefix(item.ProductNumber, "W") || strings.HasPrefix(item.ProductNumber, "K") {
				EolArtikel = append(EolArtikel, item.Id)
			}
		}
	}

	Hersteller = removeDuplicate(Hersteller)

	a.logger.Info("Successfully Sort Products")
	return NeueArtikel, AlteArtikel, EolArtikel, Hersteller, nil
}

func removeDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
