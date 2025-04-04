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
