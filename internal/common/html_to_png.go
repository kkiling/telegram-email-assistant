package common

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func saveHtml(textHtml string, dir string) error {
	filePath := filepath.Join(dir, "index.html")
	f, err := os.Create(filePath)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(textHtml)

	if err != nil {
		return err
	}

	return nil
}

func HtmlToPng(textHtml string, dir string) (string, error) {
	imgPath := filepath.Join(dir, "index.png")
	if _, err := os.Stat(imgPath); errors.Is(err, os.ErrNotExist) == false {
		return imgPath, err
	}

	textHtml = strings.ReplaceAll(textHtml, "src=\"cid:", fmt.Sprintf("src=\"%s/", dir))
	err := saveHtml(textHtml, dir)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("python", "html2png.py", dir)
	err = cmd.Run()
	if err != nil {
		log.Println(err)
		return "", err
	}

	/*stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
	}
	err = cmd.Start()
	if err != nil {
		log.Println(err)
	}
	go copyOutput(stdout)
	go copyOutput(stderr)

	err = cmd.Wait()
	if err != nil {
		log.Println(err)
	}*/

	return imgPath, nil
}

func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
