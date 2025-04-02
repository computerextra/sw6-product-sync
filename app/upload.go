package app

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
)

func (a App) UploadImages() error {
	entries, err := os.ReadDir(ImageFolder)
	if err != nil {
		return err
	}
	c, err := ftp.Dial(fmt.Sprintf("%s:21", a.env.FTP_HOST), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return fmt.Errorf("failed to dial server: %s", err.Error())
	}
	err = c.Login(a.env.FTP_USER, a.env.FTP_PASSWORD)
	if err != nil {
		return fmt.Errorf("failed to login to server: %s", err.Error())
	}
	err = c.ChangeDir(a.env.FTP_PATH)
	if err != nil {
		return fmt.Errorf("failed to change directory: %s", err.Error())
	}
	for _, e := range entries {
		f, err := os.Open(path.Join(ImageFolder, e.Name()))
		if err != nil {
			return fmt.Errorf("failed open file: %s", err.Error())
		}
		err = c.Stor(filepath.Base(f.Name()), f)
		if err != nil {
			return fmt.Errorf("failed to STOR file \"%s\": %s", f.Name(), err.Error())
		}
	}

	if err := c.Quit(); err != nil {
		return fmt.Errorf("failed to exit server: %s", err.Error())
	}

	a.logger.Info("Image Files Uploaded", slog.Any("files", len(entries)))
	return nil
}
