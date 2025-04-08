package app

import (
	"log/slog"
	"time"
)

func (a App) Delete_Products() {
	prods, err := a.getAllProducts()
	if err != nil {
		panic(err)
	}
	var payload []string
	count := 0

	for _, prod := range prods.Data {
		if count >= MAXUPLOADS {
			a.send_delete_payload(payload)
			count = 0
			payload = []string{}
			time.Sleep(10 * time.Second)
		}
		count = count + 1
		payload = append(payload, prod.Id)
	}
	if len(payload) > 0 {
		a.send_delete_payload(payload)
	}
	a.logger.Info("Successfully deleted Products", slog.Any("anzahl", len(prods.Data)))
}

func (a App) send_delete_payload(payload []string) {
	_, err := a.client.Repository.Product.Delete(a.ctx, payload)
	if err != nil {
		panic(err)
	}
}
