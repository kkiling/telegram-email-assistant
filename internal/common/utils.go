package common

import (
	"fmt"
	"github.com/flytam/filenamify"
	"os"
	"path/filepath"
)

func CreateFolderForEmailUser(fileStorageDir string, emailUser string, msgUID uint32) (string, error) {
	emailUser, err := filenamify.Filenamify(emailUser, filenamify.Options{})
	if err != nil {
		return "", fmt.Errorf("error convert %s to valid filename with error: %w", emailUser, err)
	}

	newPath := filepath.Join(fileStorageDir, emailUser, fmt.Sprint(msgUID))
	err = os.MkdirAll(newPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error create directory %s with error %w", newPath, err)
	}

	return newPath, nil
}
