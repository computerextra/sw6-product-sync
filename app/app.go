package app

import (
	"context"
	"log/slog"
	"strings"
	"sync"

	"github.com/computerextra/sw6-product-sync/config"
	"github.com/computerextra/sw6-product-sync/env"
	"github.com/computerextra/sw6-product-sync/kosatec"
	"github.com/computerextra/sw6-product-sync/shopware"
	"github.com/computerextra/sw6-product-sync/wortmann"
	sdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
	"github.com/gofrs/uuid/v5"
)

type App struct {
	ctx    context.Context
	client *sdk.Client
	env    *env.Env
	logger *slog.Logger
	config *config.Config
}

func New(logger *slog.Logger) (*App, error) {

	env, err := env.Get()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	conf, err := config.New()
	if err != nil {
		return nil, err
	}

	creds := sdk.NewPasswordCredentials(env.SW6_ADMIN_USERNAME, env.SW6_ADMIN_PASSWORD, []string{})
	client, err := sdk.NewApiClient(ctx, env.BASE_URL, creds, nil)
	if err != nil {
		return nil, err
	}

	return &App{
		ctx:    ctx,
		client: client,
		logger: logger,
		env:    env,
		config: conf,
	}, nil
}

func (a App) Uuid(name string) (string, error) {
	namespace, err := uuid.FromString(a.env.MY_NAMESPACE)
	if err != nil {
		return "", nil
	}
	u1 := uuid.NewV5(namespace, name)
	return strings.ReplaceAll(u1.String(), "-", ""), nil
}

func (a App) getAllProducts() (*sdk.ProductCollection, error) {
	apiContext := sdk.NewApiContext(a.ctx)
	criteria := sdk.Criteria{}
	criteria.Filter = []sdk.CriteriaFilter{{Type: "equals", Field: "parentId", Value: nil}}
	criteria.Limit = 500
	collection, _, err := a.client.Repository.Product.SearchAll(apiContext, criteria)
	if err != nil {
		return nil, err
	}
	a.logger.Info("Successfully downloaded shop products", slog.Any("items", len(collection.Data)))
	return collection, nil
}

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

func (a App) SynHersteller(Hersteller []string) error {
	apiContext := sdk.NewApiContext(a.ctx)
	var entity []sdk.ProductManufacturer
	for _, x := range Hersteller {
		id, err := a.Uuid(x)
		if err != nil {
			return err
		}
		entity = append(entity, sdk.ProductManufacturer{
			Id:   id,
			Name: x,
		})
	}
	_, err := a.client.Repository.ProductManufacturer.Upsert(apiContext, entity)
	if err != nil {
		return err
	}
	a.logger.Info("successfully synchronized manufacturer")

	return nil
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
