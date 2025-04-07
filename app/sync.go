package app

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/computerextra/sw6-product-sync/shopware"
	sdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
)

const MAXUPLOADS = 1000

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

func (a App) SyncCategories(neue, alte []shopware.Artikel) error {
	err := a.syncMainMenu()
	if err != nil {
		a.logger.Error("failed to sync main menu", slog.Any("error", err))
		return err
	}

	err = a.syncProductCategories(neue, alte)
	if err != nil {
		a.logger.Error("failed to sync product categories", slog.Any("error", err))
		return err
	}

	a.logger.Info("successfully synchronized categories")
	return nil
}

func (a App) syncMainMenu() error {
	apiContext := sdk.NewApiContext(a.ctx)
	var entity []sdk.Category

	// Hauptmenü anlegen nach Config
	for _, category := range a.config.Kategorien {
		cat, err := a.categoryHelper(category.Name, a.env.MAIN_CATEGORY_ID)
		if err != nil {
			a.logger.Error("failed to create category", slog.Any("error", err), slog.Any("category", category.Name))
			continue
		}
		entity = append(entity, cat)
		if len(category.Unterkategorien) > 0 {
			for _, sub := range category.Unterkategorien {
				catId, err := a.Uuid(category.Name)
				if err != nil {
					a.logger.Error("failed to create uuid", slog.Any("error", err), slog.Any("category", category.Name))
					continue
				}
				subcat, err := a.categoryHelper(sub, catId)
				if err != nil {
					a.logger.Error("failed to create category", slog.Any("error", err), slog.Any("category", category.Name))
					continue
				}
				entity = append(entity, subcat)
			}
		}
	}
	if len(entity) > 0 {
		_, err := a.client.Repository.Category.Upsert(apiContext, entity)
		if err != nil {
			a.logger.Error("failed to upsert categories", slog.Any("error", err))
			return err
		}
	}
	return nil
}

func (a App) syncProductCategories(neue, alte []shopware.Artikel) error {
	alle := append(neue, alte...)
	apiContext := sdk.NewApiContext(a.ctx)
	var entity []sdk.Category

	for _, item := range alle {
		ignored := false
		for _, ignore := range a.config.Ignore.Kategorien {
			if item.Kategorie1 == ignore {
				ignored = true
				continue
			}
		}
		if !ignored {
			for _, override := range a.config.Override {
				if item.Kategorie1 == override.AlterName && override.Index == 1 {
					item.Kategorie1 = override.NeuerName
				}
				if len(item.Kategorie2) > 0 && item.Kategorie2 == override.AlterName && override.Index == 2 {
					item.Kategorie2 = override.NeuerName
				}
				if len(item.Kategorie3) > 0 && item.Kategorie3 == override.AlterName && override.Index == 3 {
					item.Kategorie3 = override.NeuerName
				}
				if len(item.Kategorie4) > 0 && item.Kategorie4 == override.AlterName && override.Index == 4 {
					item.Kategorie4 = override.NeuerName
				}
				if len(item.Kategorie5) > 0 && item.Kategorie5 == override.AlterName && override.Index == 5 {
					item.Kategorie5 = override.NeuerName
				}
				if len(item.Kategorie6) > 0 && item.Kategorie6 == override.AlterName && override.Index == 6 {
					item.Kategorie6 = override.NeuerName
				}
			}
			item = check_dulicate_categories(item)
			parentId := ""
			for _, category := range a.config.Kategorien {
				if item.Kategorie1 == category.Name {
					parentId = a.env.MAIN_CATEGORY_ID
				}
				for _, sub := range category.Unterkategorien {
					if item.Kategorie1 == sub {
						id, err := a.Uuid(category.Name)
						if err != nil {
							a.logger.Error("failed to create uuid", slog.Any("error", err))
							continue
						}
						parentId = id
					}
				}
			}
			if len(parentId) == 0 {
				parentId = a.env.MAIN_CATEGORY_ID
			}
			cat, err := a.categoryHelper(item.Kategorie1, parentId)
			if err != nil {
				a.logger.Error("failed to create category", slog.Any("error", err), slog.Any("category", item.Kategorie1))
				continue
			}
			entity = append(entity, cat)

			cat, err = a.generateChildren(item.Kategorie2, item.Kategorie1)
			if err == nil {
				entity = append(entity, cat)
			}
			cat, err = a.generateChildren(item.Kategorie3, item.Kategorie2)
			if err == nil {
				entity = append(entity, cat)
			}
			cat, err = a.generateChildren(item.Kategorie4, item.Kategorie3)
			if err == nil {
				entity = append(entity, cat)
			}
			cat, err = a.generateChildren(item.Kategorie5, item.Kategorie4)
			if err == nil {
				entity = append(entity, cat)
			}
			cat, err = a.generateChildren(item.Kategorie6, item.Kategorie5)
			if err == nil {
				entity = append(entity, cat)
			}
		}
	}
	if len(entity) > 0 {
		_, err := a.client.Repository.Category.Upsert(apiContext, entity)
		if err != nil {
			a.logger.Error("failed to upsert categories", slog.Any("error", err), slog.Any("data", entity))
			return err
		}
	}
	return nil
}

func check_dulicate_categories(item shopware.Artikel) shopware.Artikel {
	strings := [6]string{item.Kategorie1, item.Kategorie2, item.Kategorie3, item.Kategorie4, item.Kategorie5, item.Kategorie6}

	res := [6]string{}

	for idx, item := range strings {
		if len(item) == 0 {
			continue
		}
		dupe := false
		idxFound := [2]int{9, 9}
		for idx2, c := range strings {
			if idx != idx2 && item == c {
				dupe = true
				idxFound[0] = idx
				idxFound[1] = idx2
				break
			}
		}
		if dupe {
			if idxFound[0] < idxFound[1] {
				strings[idxFound[1]] = ""
			} else {
				strings[idxFound[0]] = ""
			}
		}
		res[idx] = strings[idx]
	}

	item.Kategorie1 = res[0]
	item.Kategorie2 = res[1]
	item.Kategorie3 = res[2]
	item.Kategorie4 = res[3]
	item.Kategorie5 = res[4]
	item.Kategorie6 = res[5]

	if len(item.Kategorie5) == 0 {
		item.Kategorie6 = ""
	}
	if len(item.Kategorie4) == 0 {
		item.Kategorie6 = ""
		item.Kategorie5 = ""
	}
	if len(item.Kategorie3) == 0 {
		item.Kategorie6 = ""
		item.Kategorie5 = ""
		item.Kategorie3 = ""
	}
	if len(item.Kategorie2) == 0 {
		item.Kategorie6 = ""
		item.Kategorie5 = ""
		item.Kategorie4 = ""
		item.Kategorie3 = ""
	}
	return item
}

func (a App) categoryHelper(name string, parentId string) (sdk.Category, error) {
	id, err := a.Uuid(name)
	if err != nil {
		a.logger.Error("failed to create uuid", slog.Any("error", err))
		return sdk.Category{}, err
	}
	return sdk.Category{
		Id:                    id,
		Active:                true,
		Name:                  name,
		ParentId:              parentId,
		DisplayNestedProducts: true,
		Type:                  "page",
		ProductAssignmentType: "product",
	}, nil
}

func (a App) generateChildren(child, parent string) (sdk.Category, error) {
	if len(parent) > 0 && len(child) > 0 {
		id, err := a.Uuid(parent)
		if err != nil {
			a.logger.Error("failed to create uuid", slog.Any("error", err))
			return sdk.Category{}, err
		}
		return a.categoryHelper(child, id)
	}
	return sdk.Category{}, fmt.Errorf("parent or child is empty")
}

func (a App) CreateProducts(neu, alt []shopware.Artikel) error {
	artikel := []shopware.Artikel{}
	artikel = append(artikel, neu[:]...)
	artikel = append(artikel, alt[:]...)

	var payloads []ProductPayload

	count := 0
	for _, item := range artikel {
		if count >= MAXUPLOADS {
			apiContext := sdk.NewApiContext(a.ctx)
			_, err = a.client.Bulk.Sync(apiContext, map[string]sdk.SyncOperation{"create-product": {
				Entity:  "product",
				Action:  "upsert",
				Payload: payloads,
			}})
			if err != nil {
				a.logger.Error(
					"failed to create new products",
					slog.Any("error", err),
				)
				return err
			}
			count = 0
			payloads = []ProductPayload{}
			time.Sleep(5 * time.Second)
		}
		count = count + 1
		id, err := a.Uuid(item.Artikelnummer)
		if err != nil {
			a.logger.Error("failed to create uuid", slog.Any("error", err), slog.Any("item", item))
			continue
		}

		Kategorie := findCategory(item)

		catId, err := a.Uuid(Kategorie)
		if err != nil {
			a.logger.Error("failed to create uuid", slog.Any("error", err), slog.Any("item", item))
			continue
		}

		vkBrutto, vkNetto := a.calculate_price(item.Ek, Kategorie, 1)

		var Aktiv bool

		if item.Bestand > 0 {
			Aktiv = true
		} else {
			Aktiv = false
		}

		herstellerId, err := a.Uuid(item.Hersteller)
		if err != nil {
			a.logger.Error("failed to create uuid", slog.Any("error", err), slog.Any("item", item))
			continue
		}

		payloads = append(payloads, ProductPayload{
			Id:    id,
			TaxId: a.env.TAX_ID,
			Price: []Price{
				{
					CurrencyId:      a.env.CURRENCY_ID,
					Net:             vkNetto,
					Gross:           vkBrutto,
					Linked:          true,
					ListPrice:       nil,
					Percentage:      nil,
					RegulationPrice: nil,
					Extensions:      []any{nil},
					ApiAlias:        "price",
				},
			},
			ProductNumber: item.Artikelnummer,
			Stock:         int64(item.Bestand),
			Name:          item.Name,
			Manufacturer: ProductManufacturer{
				Id: herstellerId,
			},
			Categories: []ProductCategoryPayload{
				{
					Id: catId,
				},
			},
			ManufacturerNumber: item.HerstellerNummer,
			Visibilities: []ProductVisibility{
				{
					SalesChannelId: a.env.SALES_CHANNEL_ID,
					Visibility:     30,
				},
			},
			Description:    item.Beschreibung,
			Active:         Aktiv,
			Ean:            item.Ean,
			ShippingFree:   false,
			DeliveryTimeId: a.env.LIEFERZEIT_ID,
		},
		)
	}

	if len(payloads) > 0 {
		apiContext := sdk.NewApiContext(a.ctx)
		_, err = a.client.Bulk.Sync(apiContext, map[string]sdk.SyncOperation{"create-product": {
			Entity:  "product",
			Action:  "upsert",
			Payload: payloads,
		}})
		if err != nil {
			a.logger.Error(
				"failed to create new products",
				slog.Any("error", err),
			)
			return err
		}
	}

	a.logger.Info("successfully created items", slog.Any("items", len(artikel)))
	return nil
}

type ProductPayload struct {
	Id                 string                   `json:"id"`
	TaxId              string                   `json:"taxId"`
	Price              []Price                  `json:"price"`
	ProductNumber      string                   `json:"productNumber"`
	Stock              int64                    `json:"stock"`
	Name               string                   `json:"name"`
	Categories         []ProductCategoryPayload `json:"categories"`
	Manufacturer       ProductManufacturer      `json:"manufacturer"`
	ManufacturerNumber string                   `json:"manufacturerNumber,omitempty"`
	Visibilities       []ProductVisibility      `json:"visibilities"`
	Description        string                   `json:"description"`
	Active             bool                     `json:"active"`
	Ean                string                   `json:"ean,omitempty"`
	ShippingFree       bool                     `json:"shippingFree"`
	DeliveryTimeId     string                   `json:"deliveryTimeId"`
}

type ProductVisibility struct {
	SalesChannelId string `json:"salesChannelId"`
	Visibility     int    `json:"visibility"`
}

type ProductManufacturer struct {
	Id string `json:"id"`
}

type ProductCategoryPayload struct {
	Id string `json:"id"`
}

type Price struct {
	CurrencyId      string        `json:"currencyId"`
	Net             float64       `json:"net"`
	Gross           float64       `json:"gross"`
	Linked          bool          `json:"linked"`
	ListPrice       interface{}   `json:"listPrice,omitempty"`
	Percentage      interface{}   `json:"percentage,omitempty"`
	RegulationPrice interface{}   `json:"regulationPrice,omitempty"`
	Extensions      []interface{} `json:"extensions,omitempty"`
	ApiAlias        string        `json:"apiAlias"`
}

func findCategory(artikel shopware.Artikel) string {
	if len(artikel.Kategorie6) > 0 {
		return artikel.Kategorie6
	} else if len(artikel.Kategorie5) > 0 {
		return artikel.Kategorie5
	} else if len(artikel.Kategorie4) > 0 {
		return artikel.Kategorie4
	} else if len(artikel.Kategorie3) > 0 {
		return artikel.Kategorie3
	} else if len(artikel.Kategorie2) > 0 {
		return artikel.Kategorie2
	} else {
		return artikel.Kategorie1
	}
}

func (a App) calculate_price(ek float64, kategorie string, count int) (float64, float64) {
	aufschlag := a.config.Aufschlag.Prozentual

	for _, cA := range a.config.Kategorie {
		for _, c := range cA {
			if kategorie == c.Kategorie {
				aufschlag = c.Prozent
				break
			}
		}
	}
	var AufschlagProzent float64 = float64(aufschlag)/100 + 1
	price := ek * AufschlagProzent
	taxPercent := (a.taxRate / 100) + 1
	vk := price * taxPercent
	vk = math.Round(vk/5) * 5
	vk = vk - 0.1

	if vk < 9.9 {

		vk = 9.9
	}
	Gewinn := (vk / taxPercent) - (ek - (5 * (float64(count))))
	count = count + 1
	if Gewinn < 5 {
		return a.calculate_price(ek+5, kategorie, count)
	}

	return vk, vk / taxPercent
}
