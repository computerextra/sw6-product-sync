package app

import (
	"archive/zip"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

const KosatecFile = "kosatec.csv"
const WortmannContent = "content.csv"
const WortmannCatalog = "productcatalog.csv"
const WortmannImages = "productimages.zip"
const ImageFolder = "temp"

func (a App) Download() error {
	n, err := download_kosatec(a.env.KOSATEC_URL)
	if err != nil {
		a.logger.Error("failed to download kosatec file", slog.Any("error", err))
		return err
	}
	a.logger.Info("downloaded kosatec file", slog.Any("bytes written", n))

	n, err = download_wortmann(fmt.Sprintf("%s:21", a.env.WORTMANN_FTP_SERVER), a.env.WORTMANN_FTP_SERVER_USER, a.env.WORTMANN_FTP_SERVER_PASSWORD)
	if err != nil {
		a.logger.Error("failed to download wortmann files", slog.Any("error", err))
		return err
	}
	a.logger.Info("downloaded wortmann files", slog.Any("bytes written", n))
	err = unzip(WortmannImages, ImageFolder)
	if err != nil {
		a.logger.Error("failed to unpack images", slog.Any("error", err))
	}
	a.logger.Info("successfully unpacked image files")
	return nil
}

func download_kosatec(url string) (int64, error) {
	out, err := os.Create(KosatecFile)
	if err != nil {
		return 0, err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func download_wortmann(server string, user string, password string) (int64, error) {
	var bytesWritten int64 = 0
	c, err := ftp.Dial(server, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return bytesWritten, err
	}
	err = c.Login(user, password)
	if err != nil {
		return bytesWritten, err
	}
	n, err := download_from_ftp("content.csv", "Preisliste", WortmannContent, server, user, password)
	if err != nil {
		return bytesWritten, err
	}
	bytesWritten = bytesWritten + n
	n, err = download_from_ftp("productcatalog.csv", "Preisliste", WortmannCatalog, server, user, password)
	if err != nil {
		return bytesWritten, err
	}
	bytesWritten = bytesWritten + n
	n, err = download_from_ftp("productimages.zip", "Produktbilder", WortmannImages, server, user, password)
	if err != nil {
		return bytesWritten, err
	}
	bytesWritten = bytesWritten + n
	return bytesWritten, nil
}

func download_from_ftp(filename string, path string, dest string, server string, user string, password string) (int64, error) {
	c, err := ftp.Dial(server, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return 0, err
	}
	err = c.Login(user, password)
	if err != nil {
		return 0, err
	}

	err = c.ChangeDir(path)
	if err != nil {
		return 0, err
	}
	res, err := c.Retr(filename)
	if err != nil {
		return 0, err
	}

	file1, err := os.Create(dest)
	if err != nil {
		return 0, err
	}
	defer file1.Close()
	fileWritten, err := io.Copy(file1, res)
	if err != nil {
		return 0, err
	}
	if err := c.Quit(); err != nil {
		return 0, err
	}
	return fileWritten, nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}
			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				return err
			}
			f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
