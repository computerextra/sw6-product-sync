package app

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

func (a App) Cleanup() error {
	warn := false
	err := os.Remove(KosatecFile)
	if err != nil {
		a.logger.Warn("file not found", slog.Any("error", err))
		warn = true
	}
	err = os.Remove(WortmannCatalog)
	if err != nil {
		a.logger.Warn("file not found", slog.Any("error", err))
		warn = true
	}
	err = os.Remove(WortmannContent)
	if err != nil {
		a.logger.Warn("file not found", slog.Any("error", err))
		warn = true
	}
	err = os.Remove(WortmannImages)
	if err != nil {
		a.logger.Warn("file not found", slog.Any("error", err))
		warn = true
	}
	err = os.RemoveAll(ImageFolder)
	if err != nil {
		a.logger.Warn("file not found", slog.Any("error", err))
		warn = true
	}
	err = delete_from_ftp(a.env.FTP_PATH, a.env.FTP_HOST, a.env.FTP_USER, a.env.FTP_PASSWORD)
	if err != nil {
		a.logger.Warn("cannot delete from Server", slog.Any("error", err))
		warn = true
	}
	if warn {
		a.logger.Warn("Failed Cleanup!")
	} else {
		a.logger.Info("Successfull cleaned up")
	}
	return nil
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
