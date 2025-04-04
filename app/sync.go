package app

import sdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"

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
