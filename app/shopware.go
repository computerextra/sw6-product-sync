package app

import (
	"log/slog"

	sdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
)

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

func (a App) get_tag_rate() (float64, error) {
	apiContext := sdk.NewApiContext(a.ctx)
	criteria := sdk.Criteria{}
	criteria.Filter = []sdk.CriteriaFilter{{Type: "equals", Field: "id", Value: a.env.TAX_ID}}
	criteria.Limit = 1
	collection, _, err := a.client.Repository.Tax.SearchAll(apiContext, criteria)
	if err != nil {
		return 0, err
	}
	return collection.Data[0].TaxRate, nil
}
