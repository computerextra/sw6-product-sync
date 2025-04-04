package app

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
)

func cleanup_worker(wg *sync.WaitGroup, which string, a App) {
	defer wg.Done()

	switch which {
	case "Kosatec":
		err := os.Remove(KosatecFile)
		if err != nil {
			a.logger.Warn("file not found", slog.Any("error", err))
		}
	case "Wortmann":
		err := os.Remove(WortmannCatalog)
		if err != nil {
			a.logger.Warn("file not found", slog.Any("error", err))
		}
		err = os.Remove(WortmannContent)
		if err != nil {
			a.logger.Warn("file not found", slog.Any("error", err))
		}
		err = os.Remove(WortmannImages)
		if err != nil {
			a.logger.Warn("file not found", slog.Any("error", err))
		}
	case "ImageFolder":
		err := os.RemoveAll(ImageFolder)
		if err != nil {
			a.logger.Warn("file not found", slog.Any("error", err))
		}
	case "FTP":
		err := delete_from_ftp(a.env.FTP_PATH, a.env.FTP_HOST, a.env.FTP_USER, a.env.FTP_PASSWORD)
		if err != nil {
			a.logger.Warn("cannot delete from Server", slog.Any("error", err))
		}
	default:
		a.logger.Error("unknown cleanup type")
	}
}

func (a App) Cleanup() {
	a.logger.Info("cleanup started")
	var wg sync.WaitGroup
	wg.Add(1)
	go cleanup_worker(&wg, "Kosatec", a)
	wg.Add(1)
	go cleanup_worker(&wg, "Wortmann", a)
	wg.Add(1)
	go cleanup_worker(&wg, "ImageFolder", a)
	wg.Add(1)
	go cleanup_worker(&wg, "FTP", a)
	a.logger.Info("cleanup done")
	wg.Wait()
}

func delete_from_ftp(path, server, user, password string) error {
	c, err := ftp.Dial(fmt.Sprintf("%s:21", server), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return fmt.Errorf("failed to dial server: %s", err.Error())
	}
	err = c.Login(user, password)
	if err != nil {
		return fmt.Errorf("failed to login to server: %s", err.Error())
	}

	err = c.RemoveDirRecur(path)
	if err != nil {
		return fmt.Errorf("failed to recursivly delete directory: %s", err.Error())
	}
	pwd, err := c.CurrentDir()
	if err != nil {
		return fmt.Errorf("failed to get current dir: %s", err.Error())
	}
	err = c.ChangeDir("/")
	if err != nil {
		return fmt.Errorf("failed to change directory to %s from %s: %s", "/", pwd, err.Error())
	}
	split := strings.Split(path, "/")
	for i := range len(split) - 1 {
		err = c.ChangeDir(split[i])
		if err != nil {
			return fmt.Errorf("failed to change directory to %s from %s: %s", split[i], pwd, err.Error())
		}
	}
	err = c.MakeDir(split[len(split)-1])
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %s", split[len(split)-1], err.Error())
	}

	if err := c.Quit(); err != nil {
		return fmt.Errorf("failed to exit server: %s", err.Error())
	}
	return nil

}
