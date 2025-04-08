package app

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
	"strings"

	"github.com/computerextra/sw6-product-sync/shopware"
	sdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
	"github.com/gofrs/uuid/v5"
)

func (a App) Uuid(name string) (string, error) {
	namespace, err := uuid.FromString(a.env.MY_NAMESPACE)
	if err != nil {
		return "", nil
	}
	u1 := uuid.NewV5(namespace, name)
	return strings.ReplaceAll(u1.String(), "-", ""), nil
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

func contains(s []shopware.Artikel, e shopware.Artikel) bool {
	for _, a := range s {
		if a.Artikelnummer == e.Artikelnummer {
			return true
		}
	}
	return false
}
func remove_duplicates(arr []shopware.Artikel) []shopware.Artikel {
	list := []shopware.Artikel{}
	for _, item := range arr {
		if !contains(list, item) {
			list = append(list, item)
		}
	}

	return list
}

func (a App) check_if_ignored(item shopware.Artikel) bool {
	for _, x := range a.config.Ignore.Kategorien {
		if item.Kategorie1 == x {
			return true
		}
		if item.Kategorie2 == x {
			return true
		}
		if item.Kategorie3 == x {
			return true
		}
		if item.Kategorie4 == x {
			return true
		}
		if item.Kategorie5 == x {
			return true
		}
		if item.Kategorie6 == x {
			return true
		}
	}
	return slices.Contains(a.config.Ignore.Produkte, item.Artikelnummer)
}

func (a App) check_category(kategorie string) string {
	for _, x := range a.config.Override {
		if strings.TrimSpace(kategorie) == strings.TrimSpace(x.AlterName) {
			return x.NeuerName
		}
	}
	return kategorie
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
