package app

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/computerextra/sw6-product-sync/shopware"
	sdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
)

func (a App) CreateImages(artikel []shopware.Artikel) error {
	// TODO: Import geht nicht, Felder sind schreibgeschützt... Muss ich wahrscheinlich so machen, wie im alten!
	for _, item := range artikel {
		if len(item.Bilder) > 0 {
			coverId := ""
			ProductId, err := a.Uuid(item.Artikelnummer)
			if err != nil {
				continue
			}
			bilder := strings.Split(item.Bilder, "|")

			for idx, bild := range bilder {
				FileName := filepath.Base(bild)
				FileSuffix := strings.Split(FileName, ".")[1]

				MediaId, err := a.Uuid(bild)
				if err != nil {
					continue
				}
				ProductMediaId, err := a.Uuid(MediaId)
				if err != nil {
					continue
				}

				if idx == 0 {
					coverId = ProductMediaId
				} else {
					coverId = ""
				}
				_, err = a.client.NewRequest(a.ctx, "POST", "/api/media", ImagePayload{
					Id:            MediaId,
					MediaFolderId: a.env.MEDIA_FOLDER_ID,
				})
				if err != nil {
					a.logger.Error("failed to create media", slog.Any("error", err))
					continue
				}
				time.Sleep(20 * time.Second)
				_, err = a.client.NewRequest(a.ctx, "POST", fmt.Sprintf("/api_action/media/%s/upload", MediaId), ImageUploadPayload{
					MediaId:   MediaId,
					Url:       bild,
					Extension: FileSuffix,
					Filename:  FileName,
				})
				if err != nil {
					a.logger.Error("failed to upload media", slog.Any("error", err))
					continue
				}
				time.Sleep(20 * time.Second)
				// TODO: Link geht nicht!
				_, err = a.client.Repository.ProductMedia.Upsert(a.ctx, []sdk.ProductMedia{
					{
						Id:        ProductMediaId,
						ProductId: ProductId,
						MediaId:   MediaId,
					},
				})
				if err != nil {
					a.logger.Error("failed to link media", slog.Any("error", err))
					continue
				}
				time.Sleep(20 * time.Second)
				if len(coverId) > 0 {
					_, err := a.client.NewRequest(a.ctx, "PATCH", fmt.Sprintf("/api/product/%s", ProductId), CoverPayload{
						CoverId: coverId,
					})
					if err != nil {
						a.logger.Error("failed to link cover media", slog.Any("error", err))
						continue
					}
				}
			}
		}
		time.Sleep(20 * time.Second)
	}

	return nil
}

type ImagePayload struct {
	Id            string `json:"id"`
	MediaFolderId string `json:"mediaFolderId"`
}

type ImageUploadPayload struct {
	MediaId   string `json:"mediaId"`
	Url       string `json:"url"`
	Extension string `json:"extension"`
	Filename  string `json:"filename"`
}

type CoverPayload struct {
	CoverId string `json:"coverId"`
}
