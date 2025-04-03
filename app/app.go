package app

import (
	"context"
	"log/slog"

	"github.com/computerextra/sw6-product-sync/config"
	"github.com/computerextra/sw6-product-sync/env"
	sdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
)

type App struct {
	ctx    context.Context
	client *sdk.Client
	env    *env.Env
	logger *slog.Logger
	config *config.Config
}

func New(logger *slog.Logger, conf *config.Config) (*App, error) {

	env, err := env.Get()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	// creds := sdk.NewPasswordCredentials(env.SW6_ADMIN_USERNAME, env.SW6_ADMIN_PASSWORD, []string{})
	// client, err := sdk.NewApiClient(ctx, env.BASE_URL, creds, nil)
	// if err != nil {
	// 	return nil, err
	// }

	return &App{
		ctx: ctx,
		// client: client,
		logger: logger,
		env:    env,
		config: conf,
	}, nil
}

func (a App) GetAllProducts() (*sdk.ProductCollection, error) {
	apiContext := sdk.NewApiContext(a.ctx)
	criteria := sdk.Criteria{}
	criteria.Filter = []sdk.CriteriaFilter{{Type: "equals", Field: "parentId", Value: nil}}
	criteria.Limit = 500
	collection, _, err := a.client.Repository.Product.SearchAll(apiContext, criteria)
	if err != nil {
		return nil, err
	}
	return collection, nil
}
