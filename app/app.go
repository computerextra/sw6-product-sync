package app

import (
	"context"
	"log/slog"

	"github.com/computerextra/sw6-product-sync/config"
	"github.com/computerextra/sw6-product-sync/env"
	sdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
)

type App struct {
	ctx     sdk.ApiContext
	client  *sdk.Client
	env     *env.Env
	logger  *slog.Logger
	config  *config.Config
	taxRate float64
}

type ProductPayload struct {
	Id                 string                   `json:"id,omitempty"`
	TaxId              string                   `json:"taxId,omitempty"`
	Price              []Price                  `json:"price,omitempty"`
	ProductNumber      string                   `json:"productNumber,omitempty"`
	Stock              int64                    `json:"stock"`
	Name               string                   `json:"name,omitempty"`
	Categories         []ProductCategoryPayload `json:"categories,omitempty"`
	Manufacturer       ProductManufacturer      `json:"manufacturer,omitempty"`
	ManufacturerNumber string                   `json:"manufacturerNumber,omitempty"`
	Visibilities       []ProductVisibility      `json:"visibilities,omitempty"`
	Description        string                   `json:"description,omitempty"`
	Active             bool                     `json:"active"`
	Ean                string                   `json:"ean,omitempty"`
	ShippingFree       bool                     `json:"shippingFree"`
	DeliveryTimeId     string                   `json:"deliveryTimeId,omitempty"`
}

type ProductVisibility struct {
	SalesChannelId string `json:"salesChannelId,omitempty"`
	Visibility     int    `json:"visibility,omitempty"`
}

type ProductManufacturer struct {
	Id string `json:"id,omitempty"`
}

type ProductCategoryPayload struct {
	Id string `json:"id,omitempty"`
}

type Price struct {
	CurrencyId      string  `json:"currencyId,omitempty"`
	Net             float64 `json:"net,omitempty"`
	Gross           float64 `json:"gross,omitempty"`
	Linked          bool    `json:"linked,omitempty"`
	ListPrice       any     `json:"listPrice,omitempty"`
	Percentage      any     `json:"percentage,omitempty"`
	RegulationPrice any     `json:"regulationPrice,omitempty"`
	Extensions      []any   `json:"extensions,omitempty"`
	ApiAlias        string  `json:"apiAlias,omitempty"`
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

	creds := sdk.NewIntegrationCredentials(env.SW6_ADMIN_CLIENT_ID, env.SW6_ADMIN_CLIENT_SECRET, []string{})
	client, err := sdk.NewApiClient(ctx, env.BASE_URL, creds, nil)
	if err != nil {
		return nil, err
	}
	apiContext := sdk.NewApiContext(ctx)

	tax, err := get_tag_rate(apiContext, client, env.TAX_ID)
	if err != nil {
		return nil, err
	}

	return &App{
		ctx:     apiContext,
		client:  client,
		logger:  logger,
		env:     env,
		config:  conf,
		taxRate: tax,
	}, nil
}

func get_tag_rate(ctx sdk.ApiContext, client *sdk.Client, taxId string) (float64, error) {

	criteria := sdk.Criteria{}
	criteria.Filter = []sdk.CriteriaFilter{{Type: "equals", Field: "id", Value: taxId}}
	criteria.Limit = 1
	collection, _, err := client.Repository.Tax.SearchAll(ctx, criteria)
	if err != nil {
		return 0, err
	}
	return collection.Data[0].TaxRate, nil
}
