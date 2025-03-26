package helper

import (
	"io"
	"os"
)

func DirIsNotEmpty(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	// Check if Folder is Empty
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return false
	}
	return true
}

func InterfaceNotEmpty(value []any) bool {
	if len(value) == 0 {
		return false
	}

	return true
}
