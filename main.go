package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/computerextra/sw6-product-sync/app"
	"github.com/computerextra/sw6-product-sync/shopware"
)

const LOG = "log.txt"

func main() {
	start := time.Now()

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

	var NeueArtikel, AlteArtikel []shopware.Artikel
	var EolArtikel, Hersteller []string

	if !stop {
		NeueArtikel, AlteArtikel, EolArtikel, Hersteller, err = App.SortProducts()
		if err != nil {
			logger.Error("failed to sort products", slog.Any("error", err))
			stop = true
		}
	}
	logger.Info(
		"Produkte zum Sync",
		slog.Any("new", len(NeueArtikel)),
		slog.Any("old", len(AlteArtikel)),
		slog.Any("eol", len(EolArtikel)),
		slog.Any("hersteller", len(Hersteller)),
	)

	if !stop {
		err = App.SynHersteller(Hersteller)
		if err != nil {
			logger.Error("failed to sync manufacturer", slog.Any("error", err))
			stop = true
		}
	}

	if !stop {
		App.Cleanup()
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logger.Info(
		"used memory",
		slog.Any("Alloc", fmt.Sprintf("%v MiB", bToMb(m.Alloc))),
		slog.Any("TotalAlloc", fmt.Sprintf("%v MiB", bToMb(m.TotalAlloc))),
		slog.Any("tSys", fmt.Sprintf("%v MiB", bToMb(m.Sys))),
		slog.Any("tNumGC", m.NumGC),
	)
	elapsed := time.Since(start)
	logger.Info(
		"runtime",
		slog.Any("started", start.String()),
		slog.Any("ended", time.Now().String()),
		slog.Any("elapsed ns", fmt.Sprintf("%v ns", elapsed.Nanoseconds())),
		slog.Any("elapsed ms", fmt.Sprintf("%v ms", elapsed.Milliseconds())),
		slog.Any("elapsed s", fmt.Sprintf("%v s", elapsed.Seconds())),
		slog.Any("elapsed min", fmt.Sprintf("%v min", elapsed.Minutes())),
	)

	f.Close()

	if err := App.SendLog(LOG, stop); err != nil {
		panic(err)
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
