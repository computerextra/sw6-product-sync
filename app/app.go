package app

import (
	"context"
	"log/slog"
	"strings"

	"github.com/computerextra/sw6-product-sync/config"
	"github.com/computerextra/sw6-product-sync/env"
	"github.com/computerextra/sw6-product-sync/kosatec"
	"github.com/computerextra/sw6-product-sync/shopware"
	"github.com/computerextra/sw6-product-sync/wortmann"
	sdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
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

func (a App) SortProducts() (NeueArtikel, AlteArtikel []shopware.Artikel, EolArtikel []string, err error) {
	Kosatec, err := a.readKosatec()
	if err != nil {
		return nil, nil, nil, err
	}
	Wortmann, err := a.readWortmann()
	if err != nil {
		return nil, nil, nil, err
	}
	ShopArtikel, err := a.getAllProducts()
	if err != nil {
		return nil, nil, nil, err
	}

	for _, item := range Kosatec {
		found := false
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

	a.logger.Info("Successfully Sort Products")
	return NeueArtikel, AlteArtikel, EolArtikel, nil
}
