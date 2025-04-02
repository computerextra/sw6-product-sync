package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/computerextra/sw6-product-sync/app"
)

const LOG = "log.txt"

func main() {
	f, err := os.OpenFile(LOG, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	logger := slog.New(
		NewCopyHandler(
			slog.NewJSONHandler(os.Stdout, nil),
			slog.NewTextHandler(f, nil),
		),
	)
	var stop = false

	App, err := app.New(logger)
	if err != nil {
		logger.Error("failed to create app", slog.Any("error", err))
		stop = true
	}
	if !stop {
		err = App.Download()
		if err != nil {
			logger.Error("failed to download Files", slog.Any("error", err))
			stop = true
		}
	}

	if !stop {
		err = App.UploadImages()
		if err != nil {
			logger.Error("failed to upload images", slog.Any("error", err))
			stop = true
		}
	}

	// collection, err := App.GetAllProducts()
	// if err != nil {
	// 	logger.Error("failed to get all Products", slog.Any("error", err))
	// }
	// fmt.Printf("Anzahl der Produkte: %v", len(collection.Data))

	if !stop {
		err = App.Cleanup()
		if err != nil {
			logger.Error("failed to cleanup files", slog.Any("error", err))
		}
	}

	f.Close()

	if err := App.SendLog(LOG, stop); err != nil {
		panic(err)
	}
}
