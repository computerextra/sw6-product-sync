package app

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/computerextra/sw6-product-sync/shopware"
	sdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
)

const MAXUPLOADS = 2000

func (a App) SynHersteller(Hersteller []string) error {
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
	_, err := a.client.Repository.ProductManufacturer.Upsert(a.ctx, entity)
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
		_, err := a.client.Repository.Category.Upsert(a.ctx, entity)
		if err != nil {
			a.logger.Error("failed to upsert categories", slog.Any("error", err))
			return err
		}
	}
	return nil
}

func (a App) syncProductCategories(neue, alte []shopware.Artikel) error {
	alle := append(neue, alte...)
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
				if item.Kategorie1 == override.AlterName {
					item.Kategorie1 = override.NeuerName
				}
				if len(item.Kategorie2) > 0 && item.Kategorie2 == override.AlterName {
					item.Kategorie2 = override.NeuerName
				}
				if len(item.Kategorie3) > 0 && item.Kategorie3 == override.AlterName {
					item.Kategorie3 = override.NeuerName
				}
				if len(item.Kategorie4) > 0 && item.Kategorie4 == override.AlterName {
					item.Kategorie4 = override.NeuerName
				}
				if len(item.Kategorie5) > 0 && item.Kategorie5 == override.AlterName {
					item.Kategorie5 = override.NeuerName
				}
				if len(item.Kategorie6) > 0 && item.Kategorie6 == override.AlterName {
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
		_, err := a.client.Repository.Category.Upsert(a.ctx, entity)
		if err != nil {
			a.logger.Error("failed to upsert categories", slog.Any("error", err), slog.Any("data", entity))
			return err
		}
	}
	return nil
}

func (a App) send_product_payload(payloads []ProductPayload) error {
	_, err = a.client.Bulk.Sync(a.ctx, map[string]sdk.SyncOperation{"create-product": {
		Entity:  "product",
		Action:  "upsert",
		Payload: payloads,
	}})
	if err != nil {
		a.logger.Error(
			"failed to sync products",
			slog.Any("error", err),
		)
		return err
	}
	return nil
}

func (a App) CreateProducts(artikel []shopware.Artikel) error {

	a.logger.Info("items to be processed:", slog.Any("no. of items", len(artikel)))
	artikel = remove_duplicates(artikel)
	a.logger.Info("items to be processed after removing duplicates:", slog.Any("no. of items", len(artikel)))

	var payloads []ProductPayload

	count := 0
	for _, item := range artikel {
		if count >= MAXUPLOADS {
			if err := a.send_product_payload(payloads); err != nil {
				fmt.Println("Failed to complete Payload; Wait for 2 minutes and try again")
				time.Sleep(2 * time.Minute)
				if err := a.send_product_payload(payloads); err != nil {
					fmt.Println("Failed to complete Payload again.")
					fmt.Println("Sync every single Produkt")
					for _, load := range payloads {
						var x []ProductPayload
						x = append(x, load)
						a.send_product_payload(x)
						fmt.Println("synced Product, wait for 20 Secs")
						time.Sleep(20 * time.Second)
					}
				}

			}
			count = 0
			payloads = []ProductPayload{}
			fmt.Printf("synced %v Products, wait for 20 Secs\n", MAXUPLOADS)
			time.Sleep(20 * time.Second)
		}

		// Check if ignored
		skip := a.check_if_ignored(item)

		count = count + 1
		id, err := a.Uuid(item.Artikelnummer)
		if err != nil {
			a.logger.Error("failed to create uuid", slog.Any("error", err), slog.Any("item", item))
			continue
		}

		var Kategorie string
		Kategorie = findCategory(item)
		Kategorie = a.check_category(Kategorie)

		catId, err := a.Uuid(Kategorie)
		if err != nil {
			a.logger.Error("failed to create uuid", slog.Any("error", err), slog.Any("item", item))
			continue
		}

		vkBrutto, vkNetto := a.calculate_price(item.Ek, Kategorie, 1)

		var Aktiv bool = false
		var Bestand int64 = 0

		if item.Bestand > 0 {
			Bestand = item.Bestand
			Aktiv = true
		} else {
			Aktiv = false
		}

		herstellerId, err := a.Uuid(item.Hersteller)
		if err != nil {
			a.logger.Error("failed to create uuid", slog.Any("error", err), slog.Any("item", item))
			continue
		}
		if len(item.Ean) == 0 {
			item.Ean = "0000000000000000"
		}

		if strings.HasPrefix(item.Artikelnummer, "W") {
			if len(item.HerstellerNummer) == 0 {
				item.HerstellerNummer = strings.TrimPrefix(item.Artikelnummer, "W")
			}

			if len(item.Beschreibung) == 0 {
				item.Beschreibung = "<br>"
			}
		}

		if len(item.HerstellerNummer) == 0 {
			item.HerstellerNummer = "n/a"
		}

		// Check everything (Nothing should be empty)

		if len(id) == 0 {
			a.logger.Warn(
				"skip product",
				slog.Any("product number", item.Artikelnummer),
				slog.Any("empty", "id"),
			)
			skip = true
		}
		if len(a.env.TAX_ID) == 0 {
			a.logger.Warn(
				"skip product",
				slog.Any("product number", item.Artikelnummer),
				slog.Any("empty", "a.env.TAX_ID"),
			)
			skip = true
		}
		if len(a.env.CURRENCY_ID) == 0 {
			a.logger.Warn(
				"skip product",
				slog.Any("product number", item.Artikelnummer),
				slog.Any("empty", "a.env.CURRENCY_ID"),
			)
			skip = true
		}
		if vkNetto < 0 {
			a.logger.Warn(
				"skip product",
				slog.Any("product number", item.Artikelnummer),
				slog.Any("empty", "vkNetto"),
			)
			skip = true
		}
		if vkBrutto < 0 {
			a.logger.Warn(
				"skip product",
				slog.Any("product number", item.Artikelnummer),
				slog.Any("empty", "vkBrutto"),
			)
			skip = true
		}
		if len(item.Artikelnummer) == 0 {
			a.logger.Warn(
				"skip product",
				slog.Any("product number", item.Artikelnummer),
				slog.Any("empty", "item.Artikelnummer"),
			)
			skip = true
		}
		if item.Bestand < 0 {
			a.logger.Warn(
				"skip product",
				slog.Any("product number", item.Artikelnummer),
				slog.Any("empty", "item.Bestand"),
			)
			skip = true
		}
		if len(item.Name) == 0 {
			a.logger.Warn(
				"skip product",
				slog.Any("product number", item.Artikelnummer),
				slog.Any("empty", "item.Name"),
			)
			skip = true
		}
		if len(herstellerId) == 0 {
			a.logger.Warn(
				"skip product",
				slog.Any("product number", item.Artikelnummer),
				slog.Any("empty", "herstellerId"),
			)
			skip = true
		}
		if len(catId) == 0 {
			a.logger.Warn(
				"skip product",
				slog.Any("product number", item.Artikelnummer),
				slog.Any("empty", "catId"),
			)
			skip = true
		}
		if len(item.HerstellerNummer) == 0 {
			a.logger.Warn(
				"skip product",
				slog.Any("product number", item.Artikelnummer),
				slog.Any("empty", "item.HerstellerNummer"),
			)
			skip = true
		}
		if len(a.env.SALES_CHANNEL_ID) == 0 {
			a.logger.Warn(
				"skip product",
				slog.Any("product number", item.Artikelnummer),
				slog.Any("empty", "a.env.SALES_CHANNEL_ID"),
			)
			skip = true
		}
		if len(item.Beschreibung) == 0 {
			a.logger.Warn(
				"skip product",
				slog.Any("product number", item.Artikelnummer),
				slog.Any("empty", "item.Beschreibung"),
			)
			skip = true
		}
		if len(item.Ean) == 0 {
			a.logger.Warn(
				"skip product",
				slog.Any("product number", item.Artikelnummer),
				slog.Any("empty", "item.Ean"),
			)
			skip = true
		}
		if len(a.env.LIEFERZEIT_ID) == 0 {
			a.logger.Warn(
				"skip product",
				slog.Any("product number", item.Artikelnummer),
				slog.Any("empty", "a.env.LIEFERZEIT_ID"),
			)
			skip = true
		}
		if !skip {
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
				Stock:         Bestand,
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
	}

	if len(payloads) > 0 {
		if err := a.send_product_payload(payloads); err != nil {
			fmt.Println("Failed to complete Payload; Wait for 2 minutes and try again")
			time.Sleep(2 * time.Minute)
			if err := a.send_product_payload(payloads); err != nil {
				fmt.Println("Failed to complete Payload again.")
				fmt.Println("Sync every single Produkt")
				for _, load := range payloads {
					var x []ProductPayload
					x = append(x, load)
					a.send_product_payload(x)
					fmt.Println("synced Product, wait for 20 Secs")
					time.Sleep(20 * time.Second)
				}
			}
		}
	}

	a.logger.Info("successfully created items", slog.Any("items", len(artikel)))
	return nil
}

func (a App) Delete_Eol(artikel []string) error {
	var ids []string

	for _, item := range artikel {
		id, err := a.Uuid(item)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	if len(ids) > 0 {
		_, err := a.client.Repository.Product.Delete(a.ctx, ids)
		if err != nil {
			return err
		}
	}
	a.logger.Info("Successfully deleted eol products", slog.Any("items", len(artikel)))
	return nil
}
