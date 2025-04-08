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
	for _, item := range artikel {
		if len(item.Bilder) > 0 {
			coverId := ""
			ProductId, err := a.Uuid(item.Artikelnummer)
			if err != nil {
				continue
			}
			bilder := strings.Split(item.Bilder, "|")
			mediaEntity := []sdk.Media{}
			productMediaEntity := []sdk.ProductMedia{}

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
				}

				mediaEntity = append(mediaEntity, sdk.Media{
					Id:            MediaId,
					MediaFolderId: a.env.MEDIA_FOLDER_ID,
					Url:           bild,
					FileExtension: FileSuffix,
				})
				productMediaEntity = append(productMediaEntity, sdk.ProductMedia{
					Id:        ProductMediaId,
					ProductId: ProductId,
					MediaId:   MediaId,
				})
			}

			// Create new Media

			_, err = a.client.Repository.Media.Upsert(a.ctx, mediaEntity)
			if err != nil {
				a.logger.Error("failed to create media", slog.Any("error", err), slog.Any("item", item.Artikelnummer))
				fmt.Println("Trying to create each media")
				for _, load := range mediaEntity {
					_, err := a.client.Repository.Media.Upsert(a.ctx, []sdk.Media{
						load,
					})
					if err != nil {
						a.logger.Error("fialed to create Media", slog.Any("error", err), slog.Any("payload", load))
						continue
					}
				}
			}
			time.Sleep(20 * time.Second)
			_, err = a.client.Repository.ProductMedia.Upsert(a.ctx, productMediaEntity)
			if err != nil {
				a.logger.Error("failed to upload media", slog.Any("error", err), slog.Any("item", item.Artikelnummer))
				fmt.Println("Trying to upload media")
				for _, load := range productMediaEntity {
					_, err := a.client.Repository.ProductMedia.Upsert(a.ctx, []sdk.ProductMedia{
						load,
					})
					if err != nil {
						a.logger.Error("fialed to upload Media", slog.Any("error", err), slog.Any("payload", load))
						a.client.Repository.Media.Delete(a.ctx, []string{load.MediaId})
						continue
					}
				}
			}
			if len(coverId) != 0 {
				time.Sleep(20 * time.Second)
				_, err := a.client.Repository.Product.Upsert(a.ctx, []sdk.Product{
					{
						Id:      ProductId,
						CoverId: coverId,
					},
				})
				if err != nil {
					a.logger.Error("failed to create cover", slog.Any("error", err), slog.Any("item", item.Artikelnummer))
				}
			}
		}
		time.Sleep(20 * time.Second)
	}

	return nil
}
