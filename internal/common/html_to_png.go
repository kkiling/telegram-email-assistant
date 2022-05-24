package common

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func saveHtml(textHtml string, dir string) error {
	filePath := filepath.Join(dir, "index.html")
	f, err := os.Create(filePath)

	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Errorf("error close %s", filePath)
		}
	}(f)

	_, err = f.WriteString(textHtml)

	if err != nil {
		return err
	}

	return nil
}

func HtmlToPng(textHtml string, dir string) (string, error) {
	imgPath := filepath.Join(dir, "index.png")
	if _, err := os.Stat(imgPath); !errors.Is(err, os.ErrNotExist) {
		return imgPath, err
	}

	textHtml2 := strings.ReplaceAll(textHtml, "src=\"cid:", fmt.Sprintf("src=\"%s/", dir))
	textHtml2 = strings.ReplaceAll(textHtml2, "src=cid:", fmt.Sprintf("src=%s/", dir))

	err := saveHtml(textHtml2, dir)
	if err != nil {
		return "", err
	}

	haveCid := "false"
	if textHtml != textHtml2 {
		haveCid = "true"
	}

	err = command("python3", "html2png.py", dir, haveCid)
	if err != nil {
		return "", fmt.Errorf("error run command: %w", err)
	}

	if _, err := os.Stat(imgPath); errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	return imgPath, nil
}
