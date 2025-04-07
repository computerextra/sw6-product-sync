package app

import (
	"context"
	"log/slog"

	"github.com/computerextra/sw6-product-sync/config"
	"github.com/computerextra/sw6-product-sync/env"
	sdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
)

type App struct {
	ctx     context.Context
	client  *sdk.Client
	env     *env.Env
	logger  *slog.Logger
	config  *config.Config
	taxRate float64
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

	tax, err := get_tag_rate(ctx, client, env.TAX_ID)
	if err != nil {
		return nil, err
	}

	return &App{
		ctx:     ctx,
		client:  client,
		logger:  logger,
		env:     env,
		config:  conf,
		taxRate: tax,
	}, nil
}

func get_tag_rate(ctx context.Context, client *sdk.Client, taxId string) (float64, error) {
	apiContext := sdk.NewApiContext(ctx)
	criteria := sdk.Criteria{}
	criteria.Filter = []sdk.CriteriaFilter{{Type: "equals", Field: "id", Value: taxId}}
	criteria.Limit = 1
	collection, _, err := client.Repository.Tax.SearchAll(apiContext, criteria)
	if err != nil {
		return 0, err
	}
	return collection.Data[0].TaxRate, nil
}
