package main

import (
	"flag"
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
	help := flag.Bool("h", false, "Show Help")
	endless := flag.Bool("endless", false, "Run Programm Endlessly")
	deleteProducts := flag.Bool("delete-products", false, "Deletes all Products")
	timer := flag.Int("wait", 0, "Time to wait between iterations")

	flag.Parse()

	if *deleteProducts {
		fmt.Println("Running in deletion Mode")
		fmt.Println("Deleting Products")
		delete_products()
	} else if *help {
		fmt.Println("Help:")
		fmt.Println("-h : Show this help Window")
		fmt.Println("-endless : Run the Programm endlessly")
		fmt.Println("-wait : time to wait between runs")
		fmt.Println("------")
		fmt.Println("Example:")
		fmt.Println("./sw6-product-sync.exe -endless -wait=2")
		fmt.Println("Programm runs, wait for 2 hours and runs again, infinite")
	} else if !*help && !*endless && *timer > 0 {
		fmt.Println("Invalid Argument Chain, -wait can only be used in conjunction with -endless")
	} else if *endless && *timer == 0 {
		fmt.Println("Invalid Argument Chain, -endless can only be used in conjunction with -wait")
	} else if *endless && *timer > 0 {
		fmt.Println("Starting Programm in endless Mode: Timer: ", *timer)
		for {
			run_program()
			x := *timer
			time.Sleep(time.Duration(x) * time.Hour)
		}
	} else {
		fmt.Println("Running without Args, running programm once")
		run_program()
	}

}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func delete_products() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	App, err := app.New(logger)
	if err != nil {
		panic(err)
	}
	App.Delete_Products()
}

func run_program() {
	start := time.Now()

	f, err := os.OpenFile(LOG, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	logger := slog.New(
		NewCopyHandler(
			slog.NewJSONHandler(os.Stdout, nil),
			slog.NewJSONHandler(f, nil),
		),
	)

	var stop = false

	App, err := app.New(logger)
	if err != nil {
		logger.Error("failed to create app", slog.Any("error", err))
		stop = true
	}

	// if !stop {
	// 	err = App.Download()
	// 	if err != nil {
	// 		logger.Error("failed to download Files", slog.Any("error", err))
	// 		stop = true
	// 	}
	// }

	// if !stop {
	// 	err = App.UploadImages()
	// 	if err != nil {
	// 		logger.Error("failed to upload images", slog.Any("error", err))
	// 		stop = true
	// 	}
	// }

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

	// if !stop {
	// 	err = App.SyncCategories(NeueArtikel, AlteArtikel)
	// 	if err != nil {
	// 		logger.Error("failed to sync categories", slog.Any("error", err))
	// 		stop = true
	// 	}
	// }

	if !stop {
		err = App.CreateProducts(AlteArtikel)
		if err != nil {
			logger.Error("failed to sync Products", slog.Any("error", err))
			stop = true
		}
	}

	if !stop {
		err = App.CreateProducts(NeueArtikel)
		if err != nil {
			logger.Error("failed to sync Products", slog.Any("error", err))
			stop = true
		}
	}

	if !stop {
		err = App.CreateImages(NeueArtikel)
		if err != nil {
			logger.Error("failed to create images", slog.Any("error", err))
			stop = true
		}
	}

	if !stop {
		err = App.Delete_Eol(EolArtikel)
		if err != nil {
			logger.Error("failed to delete eol products", slog.Any("error", err))
			stop = true
		}
	}

	if !stop {
		// App.Cleanup()
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

	// if err := App.SendLog(LOG, stop); err != nil {
	// 	panic(err)
	// }
}
