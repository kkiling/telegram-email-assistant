package common

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/flytam/filenamify"
)

func CreateFolderForEmail(fileStorageDir string, to string, msgUID int64) (string, error) {
	emailUser, err := filenamify.Filenamify(to, filenamify.Options{})
	if err != nil {
		return "", fmt.Errorf("error convert %s to valid filename with error: %w", emailUser, err)
	}

	newPath := filepath.Join(fileStorageDir, to, fmt.Sprint(msgUID))
	err = os.MkdirAll(newPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error create directory %s with error %w", newPath, err)
	}

	return newPath, nil
}
