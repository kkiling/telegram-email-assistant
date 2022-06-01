package common

import (
	"bufio"
	"io"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func printStdout(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		logrus.Debugf("stdout: %s", scanner.Text())
	}
}

func printStderr(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		logrus.Warnf("stderr: %s", scanner.Text())
	}
}

func command(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	go printStdout(stdout)
	go printStderr(stderr)

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
